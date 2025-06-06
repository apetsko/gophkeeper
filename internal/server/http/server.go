// Package http provides the HTTP server implementation for the GophKeeper service.
//
// It sets up an HTTP server with gRPC-Gateway for proxying RESTful requests to the gRPC backend,
// configures CORS middleware, and manages TLS support. The package handles server startup,
// graceful shutdown, and integration with application configuration and logging.
package http

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apetsko/gophkeeper/config"
	"github.com/apetsko/gophkeeper/pkg/logging"
	pb "github.com/apetsko/gophkeeper/protogen/api/proto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// RunHTTP starts the HTTP server with gRPC-Gateway and CORS support.
//
// This function configures the HTTP server to proxy requests to the gRPC backend using gRPC-Gateway,
// sets up CORS middleware, and manages TLS if enabled. It listens on the configured address and
// gracefully shuts down the server on context cancellation.
//
// Parameters:
//   - ctx: The context for server lifecycle management.
//   - cfg: Pointer to the application configuration.
//   - log: Logger instance for server logging.
//
// Returns:
//   - *http.Server: The running HTTP server instance.
//   - error: An error if the server fails to start.
func RunHTTP(ctx context.Context, cfg *config.Config, log *logging.Logger) (*http.Server, error) {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.New(nil)
			if auth := req.Header.Get("jwt"); auth != "" {
				md.Set("jwt", auth)
			}
			return md
		}),
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	var opts []grpc.DialOption

	if cfg.TLSConfig.EnableHTTPS {
		// Загрузить сертификат CA (корневой сертификат), которым подписан сервер gRPC,
		// чтобы grpc-gateway доверял этому соединению
		certPool := x509.NewCertPool()
		caCert, err := os.ReadFile(cfg.TLSConfig.CertPath)
		if err != nil {
			log.Fatalf("failed to read CA cert: %v", err)
		}
		if !certPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("failed to append CA cert")
		}
		creds := credentials.NewClientTLSFromCert(certPool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err := pb.RegisterGophKeeperHandlerFromEndpoint(ctx, mux, cfg.GRPCAddress, opts); err != nil {
		log.Fatalf("failed to register gRPC-Gateway: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		<-ctx.Done()
		five := 5 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), five)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		log.Info(fmt.Sprintf("Starting HTTP server at %s, TLS: %t", srv.Addr, cfg.TLSConfig.EnableHTTPS))
		if cfg.TLSConfig.EnableHTTPS {
			return srv.ListenAndServeTLS(cfg.TLSConfig.CertPath, cfg.TLSConfig.KeyPath)
		}
		return srv.ListenAndServe()
	})

	go func() {
		if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error: " + err.Error())
		}
	}()

	return srv, nil
}
