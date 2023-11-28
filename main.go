package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/kataras/iris/v12"
)

// DynamicCertificates is a structure to hold dynamic TLS certificates.
type DynamicCertificates struct {
	sync.RWMutex
	Certificates map[string]tls.Certificate
}

func main() {
	initVariable()
	app := iris.New()
	appHTTP := iris.New()

	// Initialize the dynamic TLS certificates map.
	dynamicCertificates := &DynamicCertificates{
		Certificates: make(map[string]tls.Certificate),
	}

	// Set up a custom TLS configuration based on the requested domain.
	tlsConfig := &http.Server{
		Addr:    ":" + httpsPort,
		Handler: app,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				domain := info.ServerName
				dynamicCertificates.RLock()
				cert, exists := dynamicCertificates.Certificates[domain]
				dynamicCertificates.RUnlock()
				if !exists {
					return nil, fmt.Errorf("no TLS certificate found for domain %s", domain)
				}
				return &cert, nil
			},
		},
	}

	// Define a route to update TLS certificates dynamically.
	app.Post("/updateCertificates", func(ctx iris.Context) {
		var requestData struct {
			Domain string `json:"domain"`
			Cert   string `json:"cert"`
			Key    string `json:"key"`
		}

		if err := ctx.ReadJSON(&requestData); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid JSON request"})
			return
		}

		// Update the TLS certificates dynamically.
		dynamicCertificates.Lock()
		dynamicCertificates.Certificates[requestData.Domain], _ = tls.X509KeyPair([]byte(requestData.Cert), []byte(requestData.Key))
		dynamicCertificates.Unlock()

		ctx.JSON(iris.Map{"message": "TLS certificates updated successfully"})
	})

	// HTTP Biasa
	go func() {
		appHTTP.Any("/{any:path}", renderMain)
		err := appHTTP.Run(iris.Addr(":" + httpPort))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// HTTPS
	app.Any("/{any:path}", renderMain)
	err := app.Run(iris.Server(tlsConfig))
	if err != nil {
		log.Fatal(err)
	}

	// Wait for an interrupt signal (e.g., Ctrl+C) to gracefully shut down the servers.
	iris.RegisterOnInterrupt(func() {
		log.Println("Shutting down gracefully...")
		app.Shutdown(nil)
		appHTTP.Shutdown(nil)
		log.Println("Shutting down success...")
	})

	log.Println("Server is running on ports 80 and 443...")
}

// redirectHandler redirects HTTP requests to the HTTPS server.
func redirectHandler(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func renderMain(ctx iris.Context) {
	//Get Domain Data
	domain, ok := GetDomain(ctx)
	if !ok {
		ctx.HTML("404")
		return
	}
	switch domain.Type {
	case "reverse_proxy":
		success := reverseProxy(ctx, domain)
		if success {
			return
		}
		ctx.HTML("Tidak ada respons, pastikan server anda aktif")
		break
	default:
		ctx.HTML("404")
		return
	}
	//Reverse Proxy
}
