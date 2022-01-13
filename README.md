# cromwell-om

## Usage
```shell
NAME: 
   Cromwell-OM - OM for Cromwell Server

USAGE:
   cromwell-om [global options] command [command options] [arguments...]

COMMANDS:
   version, v  Cromwell-OM version
   status, s   Status of workflow, auto return failed info
   summary, m  Summary through metadata data (table)
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value  Url for your Cromwell Server (default: "http://127.0.0.1:8000")
   --help, -h    show help (default: false)
```
## Version
```shell
Author: wangxuehan
v1.1
```
## Config
Must be named **cromwell-om.conf**
```json
{
  "host" : "http://127.0.0.1:8000"
}
```

## Example
### status
```shell
# show workflow status, and show the failed message if workflow failed
cromwell-om 
  status \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c
# or
cromwell-om s -i c0c9471b-cf62-4215-900e-f671ac318e2c
```
![image](https://user-images.githubusercontent.com/70520563/149335064-5f616ede-6382-442d-bd9a-915384abe79d.png)

***
### summary
```shell
# show tasks info
cromwell-om 
  summary \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c
# or
cromwell-om m -i c0c9471b-cf62-4215-900e-f671ac318e2c
```
![image](https://user-images.githubusercontent.com/70520563/148772099-d3e3f34b-d592-4de7-b669-552eb3dbe158.png)
***
```shell
# show tasks info and the time resuming
cromwell-om 
  summary \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c \
  -t
```
![image](https://user-images.githubusercontent.com/70520563/148772175-6a466528-19c8-43f4-9fc7-57f82a753a2c.png)
![image](https://user-images.githubusercontent.com/70520563/149336180-f293217f-7d16-42ea-9bd2-8ed0e96e77c8.png)

***
```shell
# show failed tasks info
cromwell-om \
  summary \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c \
  -f
```
![image](https://user-images.githubusercontent.com/70520563/148772216-d4e94e27-20a7-4658-9acf-abbf9984aa66.png)
***
```shell
# show failed tasks info and the path of stderr
cromwell-om \
  summary \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c \
  -f -e
```
![image](https://user-images.githubusercontent.com/70520563/149335452-07a63156-3629-44b5-8e14-28a40817cca1.png)
```shell
# show sub workflow
cromwell-om \
  summary \
  -i c0c9471b-cf62-4215-900e-f671ac318e2c \
  -s
```
![image](https://user-images.githubusercontent.com/70520563/149335892-837c6768-502d-4722-ad36-ba599aa201ac.png)

