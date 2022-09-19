package introspection

import (
	"fmt"
	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
)

func NewIntrospectEngine() *IntrospectEngine {
	// TODO:这里可以设置默认值
	engine := &IntrospectEngine{
		// path: path,
	}
	file, err := engine.ensure() //确保引擎一定安装了
	if err != nil {
		panic(err)
	}
	engine.path = file
	return engine
}

type IntrospectEngine struct {
	path string
}

func (e *IntrospectEngine) ensure() (string, error) {
	ensureEngine := time.Now()

	dir := binaries.GlobalCacheDir()
	// 确保引擎一定下载了
	if err := binaries.FetchNative(dir); err != nil {
		return "", fmt.Errorf("could not fetch binaries: %w", err)
	}
	// check for darwin/windows/linux first
	binaryName := platform.CheckForExtension(platform.Name(), platform.Name())

	var file string
	// forceVersion saves whether a version check should be done, which should be disabled
	// when providing a custom query engine value
	// forceVersion := true
	name := "prisma-introspection-engine-"
	globalPath := path.Join(dir, binaries.EngineVersion, name+binaryName)

	logger.Debug.Printf("expecting global introspection engine `%s` ", globalPath)

	// TODO write tests for all cases
	// first, check if the query engine binary is being overridden by PRISMA_MIGRATION_ENGINE_BINARY
	prismaQueryEngineBinary := os.Getenv("PRISMA_INTROSPECTION_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.Debug.Printf("PRISMA_INTROSPECTION_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)

		if _, err := os.Stat(prismaQueryEngineBinary); err != nil {
			return "", fmt.Errorf("PRISMA_INTROSPECTION_ENGINE_BINARY was provided, but no query engine was found at %s", prismaQueryEngineBinary)
		}

		file = prismaQueryEngineBinary
		// forceVersion = false
	} else {
		if _, err := os.Stat(globalPath); err == nil {
			logger.Debug.Printf("exact query engine found in global path")
			file = globalPath
		}
	}

	if file == "" {
		// TODO log instructions on how to fix this problem
		return "", fmt.Errorf("no binary found ")
	}
	logger.Debug.Printf("using introspection engine at %s", file)
	logger.Debug.Printf("ensure introspection engine took %s", time.Since(ensureEngine))

	return file, nil
}

func (e *IntrospectEngine) Pull(schemaPath string) error {
	startParse := time.Now()
	// 可以缓存到改引擎中？
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalln("load prisma schema", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.path, "--datamodel", schemaPath)
	//cmd := exec.CommandContext(ctx, e.path, "--datamodel", schemaPath)

	pipe, err := cmd.StdinPipe() // 标准输入流
	if err != nil {
		log.Fatalln("introspect engine std in pipe", err)
		return err
		// return "", err
	}
	defer pipe.Close()
	// 构建一个json-rpc 请求参数
	req := IntrospectRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "compositeTypeDepth",
		Params: IntrospectRequestParams{
			CompositeTypeDepth: -1,
			Schema:             string(schema),
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// 入参追加到管道中
	_, err = pipe.Write(append(data, []byte("\n")...))
	if err != nil {
		// return "", err
		return err
	}
	// 开始执行
	err = cmd.Start()
	if err != nil {
		return err
	}

	var response IntrospectResponse

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("Introspect std out pipe", err)
	}
	r := bufio.NewReader(out)
	outBuf := &bytes.Buffer{}
	// {\"jsonrpc\":\"2.0\",\"result\":{\"executedSteps\":1,\"unexecutable\":[],\"warnings\":[]},\"id\":1}\n
	// 这一段的意思是，每100ms读取一次结果，直到超时或有结果
	for {
		// 等待100 ms
		time.Sleep(time.Millisecond * 100)
		b, err := r.ReadByte()
		if err != nil {
			log.Fatalln("introspect ReadByte", err)
		}
		err = outBuf.WriteByte(b)
		if err != nil {
			log.Fatalln("introspect writeByte", err)
		}

		if b == '\n' {
			// 解析响应结果
			err = json.Unmarshal(outBuf.Bytes(), &response)
			if err != nil {
				return err
			}
			if response.Error == nil {
				log.Println("introspect successful")
			}
			fmt.Print("ende ")
			break
		}
		// 如果超时了？跳出读取？
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	log.Printf("[timing] introspect took %s", time.Since(startParse))
	if response.Error != nil {
		return fmt.Errorf("introspect error: %s", response.Error.Message)
	}
	return nil
}

func (e *IntrospectEngine) Pull2(schemaPath string) error {
	startParse := time.Now()

	// 可以缓存到改引擎中？
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalln("load prisma schema", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	//cmd := exec.CommandContext(ctx, e.path, "--datamodel", schemaPath)
	cmd := exec.CommandContext(ctx, e.path)

	pipe, err := cmd.StdinPipe() // 标准输入流
	if err != nil {
		log.Fatalln("Introspect engine std in pipe", err)
		return err
		// return "", err
	}
	defer pipe.Close()
	// 构建一个json-rpc 请求参数
	req := IntrospectRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "introspect",
		Params: IntrospectRequestParams{
			CompositeTypeDepth: -1,
			Schema:             string(schema),
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		// return "", err
		return err

	}
	// 入参追加到管道中
	_, err = pipe.Write(append(data, []byte("\n")...))
	if err != nil {
		return err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("Introspect std out pipe", err)
	}

	var response IntrospectResponse

	go func() {
		r := bufio.NewReader(out)
		outBuf := &bytes.Buffer{}
		// {\"jsonrpc\":\"2.0\",\"result\":{\"executedSteps\":1,\"unexecutable\":[],\"warnings\":[]},\"id\":1}\n
		// 这一段的意思是，每100ms读取一次结果，直到超时或有结果
		for {
			// 等待100 ms
			time.Sleep(time.Millisecond * 100)
			b, err := r.ReadByte()
			if err != nil {
				log.Fatalln("Introspect ReadByte", err)
			}
			err = outBuf.WriteByte(b)
			if err != nil {
				log.Fatalln("Introspect writeByte", err)
			}

			if b == '\n' {
				cancel() //终止进程
				// 解析响应结果
				err = json.Unmarshal(outBuf.Bytes(), &response)
				if err != nil {
					log.Fatalln("Introspect unmarshal response", err)
					return
				}
				fmt.Print("read complete ")
				return
			}
		}
	}()
	// 阻塞运行
	err = cmd.Run()
	log.Printf("[timing] Introspect2 took %s", time.Since(startParse))

	if err != nil && ctx.Err() == nil {
		log.Println("Introspect engine run", err)
		return fmt.Errorf("introspect error: %s", err)
	}
	if response.Error != nil {
		return fmt.Errorf("introspect error: %s", response.Error.Message)
	}
	return nil
}
