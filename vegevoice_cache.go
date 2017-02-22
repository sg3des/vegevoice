package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

var dirCache string

var cacheFavicons map[string]*gdkpixbuf.Pixbuf

func SetCacheDir(dir string) {
	cacheFavicons = make(map[string]*gdkpixbuf.Pixbuf)
	dirCache = path.Join(dir, "cache")
	err := os.MkdirAll(dirCache, 0755)
	if err != nil {
		log.Println(err)
	}
}

func GetFavicon(domain string, iconuri string) *gdkpixbuf.Pixbuf {
	log.Println(domain, iconuri)
	if len(domain) == 0 {
		return nil
	}

	if len(iconuri) > 0 {
		if icon, ok := cacheFavicons[domain]; ok {
			log.Println("return from cache", domain)
			return icon
		}
	}

	iconpath := getFaviconPath(domain)
	if _, err := os.Stat(iconpath); err == nil {
		//get from disk
		icon := resizeFavicon(iconpath)
		if icon != nil {
			log.Println("return from disk", domain)
			cacheFavicons[domain] = icon
			return icon
		}
	}

	if len(iconuri) > 0 {
		//download favicon
		icon, err := downloadFavicon(domain, iconuri)
		if err != nil {
			log.Println(err)
			return nil
		}
		log.Println("download done for", domain)

		cacheFavicons[domain] = icon
		return icon
	}
	return nil
}

func getFaviconPath(domain string) string {
	return path.Join(dirCache, domain, domain+".favicon")
}

func resizeFavicon(iconpath string) *gdkpixbuf.Pixbuf {
	pix := gtk.NewImageFromFile(iconpath).GetPixbuf()
	return pix.ScaleSimple(12, 12, gdkpixbuf.INTERP_BILINEAR)
}

func downloadFavicon(domain string, uri string) (*gdkpixbuf.Pixbuf, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	iconpath := getFaviconPath(domain)
	os.MkdirAll(path.Dir(iconpath), 0755)

	iconfile, err := os.Create(iconpath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(iconfile, response.Body)
	if err != nil {
		return nil, err
	}
	iconfile.Close()

	return resizeFavicon(iconpath), nil
}
