# smoketester

![example workflow](https://github.com/smutil/smoketester/actions/workflows/build-actions.yml/badge.svg)![example workflow](https://github.com/smutil/smoketester/actions/workflows/release-actions.yml/badge.svg)

Tool to execute smoke tests. smoketester can be used in CI/CD pipeline for smoketesting after deployment to test api response and application health check.

Current Features
------------
1. Multiple targets/endpoints can be tested
2. Configurable retry option in case of test failure
3. Configurable Qualitygate
4. Option to provide multiple texts in response for matching
5. GET, POST and PUT methods are supported
6. Basic authentication supported
7. Environment variables can be defined as ${ENVNAME}


Usage
-----
 step 1. download smoketester from <a href=https://github.com/smutil/smoketester/releases>releases</a>. 
 
 step 2. create [config.yml](https://github.com/smutil/smoketester/blob/main/config.yml). If config.yaml and smoketester is not in same location, you can provide the config.yml path using --config
 
 step 3. execute the smoketester as shown below. 
 
 ```
 ./smoketester --config /path-to-config.yml
 ```

Configuration
-----

 ```
global:
  [ retry: <retry count , can be overwritten for each test > | type:int ]
  [ retryInterval: <retry interval in seconds, can be overwritten for each test> | type:int ]
  [ statusCode: <http return status code, can be overwritten for each test> | type: int ]
  [ qualitygate: <percentage of success tests as qualitygate> | type: int ]
targets:  <list of targets for testing>
  - [ url: < http/https target endpoint > | type:string ]
    [ name : < testname > | type:string ]
    [ method : < GET/POST/PUT > | type:string ]
    [ username : < username, needed for basic authentication > | type:string ]
    [ password : < password, needed for basic authentication > | type:string ]
    [ header : < ["Content-Type application/json", "Accept application/json" ,.....]  | type: array of string]
    [ data : < input datafile path for POST/PUT request > | type:string ]
    [ responseText : < ["healthy", "UP" ,.....]  | type: array of string]
    [ retry: <retry count > | type:int ]
    [ retryInterval: <retry interval in seconds> | type:int ]
    [ statusCode: <http return status code> | type: int ]
  
   ```
