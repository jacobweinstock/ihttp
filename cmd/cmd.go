package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ardanlabs/service/foundation/web"
	"github.com/jacobweinstock/ihttp/foundation/content/file"
	chttp "github.com/jacobweinstock/ihttp/foundation/content/http"
	lfile "github.com/jacobweinstock/ihttp/foundation/locate/file"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	CBackend  string
	CFilePath string
	LBackend  string
	Port      int
}

type Handler struct {
	content ContentReader
	locator Locator
}

type Locator interface {
	Open(context.Context) (io.ReadCloser, error)
	Locate(context.Context, net.HardwareAddr, io.Reader) (string, error)
	Close(ctx context.Context, rc io.ReadCloser) error
}

type ContentReader interface {
	Open(context.Context, string) (io.ReadCloser, error)
	Read(ctx context.Context, r io.Reader) ([]byte, error)
	Close(ctx context.Context, rc io.ReadCloser) error
}

func (h *Handler) autoWithCtx(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	log.Println("handling request")

	m, err := net.ParseMAC(web.Param(req, "mac"))
	if err != nil {
		log.Printf("error parsing mac: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	o, err := h.locator.Open(ctx)
	if err != nil {
		log.Printf("error opening location: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer h.locator.Close(ctx, o)
	loc, err := h.locator.Locate(ctx, m, o)
	if err != nil {
		log.Printf("error getting location: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	rc, err := h.content.Open(ctx, loc)
	if err != nil {
		log.Printf("error opening content: %v: loc: %v", err, loc)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer h.content.Close(ctx, rc)
	content, err := h.content.Read(ctx, rc)
	if err != nil {
		log.Printf("error getting content: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(content)
	return nil
}

func Run(ctx context.Context, c Config) error {

	return ultimateService(ctx, c)
}

func ultimateService(ctx context.Context, c Config) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	w := web.NewApp(shutdown)

	var cr ContentReader
	switch c.CBackend {
	case "file":
		cr = &file.Config{}
	case "http":
		cr = &chttp.Config{}
	default:
		return fmt.Errorf("unknown backend: %q", c.CBackend)
	}

	var l Locator
	switch c.LBackend {
	case "file":
		l = &lfile.Config{URI: c.CFilePath}
	default:
		return fmt.Errorf("unknown backend: %q", c.LBackend)
	}

	h := &Handler{content: cr, locator: l}
	w.Handle(http.MethodGet, "", "/:mac/auto.ipxe", h.autoWithCtx)
	s := http.Server{
		Addr:    ":" + strconv.Itoa(c.Port),
		Handler: w,
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		log.Println("shutting down")
		return s.Shutdown(ctx)
	})
	return g.Wait()
}

/*
func (h *Handler) auto(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println("handling request")

	mac := mux.Vars(req)["mac"]
	log.Println(" here mac:", mac)
	m, err := net.ParseMAC(mac)
	if err != nil {
		log.Printf("error parsing mac: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := os.Open(h.reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cfg.Close()
	l := &lfile.Config{Reader: cfg}
	loc, err := h.getLocation(ctx, l, m)
	if err != nil {
		log.Printf("error getting location: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := os.Open(loc)
	if err != nil {
		log.Printf("error opening location: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer data.Close()
	c := &file.Config{Reader: data}
	content, err := h.getContent(ctx, c)
	if err != nil {
		log.Printf("error getting content: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(content))
}

func std(ctx context.Context, file string) error {
	h := &Handler{reader: file}
	m := mux.NewRouter()
	m.HandleFunc("/{mac}/auto.ipxe", h.auto)

	s := http.Server{
		Addr:    ":8080",
		Handler: m,
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		log.Println("shutting down")
		return s.Shutdown(ctx)
	})
	return g.Wait()
}
*/
