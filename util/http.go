package util

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func PostHttpResult(ctx context.Context, url string, param interface{}, result interface{}) error {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return err
	}
	resultBytes, err := post(ctx, url, "application/json;charset=UTF-8", paramBytes)
	if err != nil {
		return err
	}
	return json.Unmarshal(resultBytes, result)
}

func PostXmlHttpResult(ctx context.Context, url string, contentType string, result interface{}, params ...interface{}) error {
	paramBytes := []byte(xml.Header)
	for _, param := range params {
		bytes, err := xml.Marshal(&param)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		paramBytes = append(paramBytes, bytes...)
	}
	resultBytes, err := post(ctx, url, contentType, paramBytes)
	if err != nil {
		return err
	}
	return xml.Unmarshal(resultBytes, result)
}

func PostXmlHttpStringResult(ctx context.Context, url string, contentType string, needHeader bool, params ...interface{}) (string, error) {
	var paramBytes []byte
	if needHeader {
		paramBytes = []byte(xml.Header)
	}
	for _, param := range params {
		bytes, err := xml.Marshal(&param)
		if err != nil {
			zap.L().Error(err.Error())
			return "", err
		}
		paramBytes = append(paramBytes, bytes...)
	}
	resultBytes, err := post(ctx, url, contentType, paramBytes)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(resultBytes), "\r\n", ""), nil
}

func PostXmlHttpBytesResult(ctx context.Context, url string, contentType string, needHeader bool, params ...interface{}) ([]byte, error) {
	var paramBytes []byte
	if needHeader {
		paramBytes = []byte(xml.Header)
	}
	for _, param := range params {
		bytes, err := xml.Marshal(&param)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		paramBytes = append(paramBytes, bytes...)
	}
	resultBytes, err := post(ctx, url, contentType, paramBytes)
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}

func GetHttpBytesResult(ctx context.Context, url string) ([]byte, error) {
	resultBytes, err := get(ctx, url)
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}

func post(ctx context.Context, url string, contentType string, param []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(param))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		span, newCtx := opentracing.StartSpanFromContext(
			ctx,
			fmt.Sprintf("http::post::%s::%s::%s", url, string(param), ""),
			opentracing.Tag{Key: "err", Value: err},
		)
		opentracing.SpanFromContext(newCtx)
		span.Finish()
		zap.L().Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span, newCtx := opentracing.StartSpanFromContext(
			ctx,
			fmt.Sprintf("http::post::%s::%s::%s", url, string(param), ""),
			opentracing.Tag{Key: "status", Value: resp.Status},
			opentracing.Tag{Key: "err", Value: err},
		)
		opentracing.SpanFromContext(newCtx)
		span.Finish()
		zap.L().Error(err.Error())
		return nil, err
	}
	span, newCtx := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("http::post::%s::%s::%s", url, string(param), string(body)),
		opentracing.Tag{Key: "status", Value: resp.Status},
		opentracing.Tag{Key: "err", Value: err},
	)
	opentracing.SpanFromContext(newCtx)
	span.Finish()
	return body, nil
}

func get(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		span, newCtx := opentracing.StartSpanFromContext(
			ctx,
			fmt.Sprintf("http::get::%s", url),
			opentracing.Tag{Key: "err", Value: err},
		)
		opentracing.SpanFromContext(newCtx)
		span.Finish()
		zap.L().Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span, newCtx := opentracing.StartSpanFromContext(
			ctx,
			fmt.Sprintf("http::get::%s", url),
			opentracing.Tag{Key: "status", Value: resp.Status},
			opentracing.Tag{Key: "err", Value: err},
		)
		opentracing.SpanFromContext(newCtx)
		span.Finish()
		zap.L().Error(err.Error())
		return nil, err
	}
	span, newCtx := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("http::get::%s", url),
		opentracing.Tag{Key: "status", Value: resp.Status},
		opentracing.Tag{Key: "err", Value: err},
	)
	opentracing.SpanFromContext(newCtx)
	span.Finish()
	return body, nil
}
