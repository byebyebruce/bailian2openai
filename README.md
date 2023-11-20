# BAILIAN2OPENAI
> [BaiLian](https://bailian.console.aliyun.com/) API to OpenAI API
 
## 目的 
目前能免费申请到Aliyun的[百炼大模型平台](https://bailian.console.aliyun.com/)服务，但已经无数开源应用都是基于OpenAI API为标准开发的   
用这个服务，可以方便接入其他以OpenAI API为标准的项目。


## Prepare
Get `Access Key、Access Secret Key、Agent Key、AppId` from [BaiLian](https://help.aliyun.com/document_detail/2587494.html)

## Run
1. Config `cp .env.exmaple .env` then edit it.
2. Run
   - Run api server `bailian2openai`
   - Run CLI `bailian2openai -cli`
   - Help `bailian2openai --help`  

## Docker
1. Build
```bash
docker build \
-t image:latest \
--build-arg VERSION=$(git describe --tags --always) \
.
```

2. Run
```bash
docker run -it --rm -p 8080:8080 \
   -e ACCESS_KEY_ID= \
   -e ACCESS_KEY_SECRET= \
   -e AGENT_KEY= \
   -e APP_ID= \
   bailu1901/bailian2openai:latest
```

## Example
- [api-server](./cmd/example/api-server)
- [chatgpt-next-web](./cmd/example/chatgpt-next-web)

## Build
`make build`
