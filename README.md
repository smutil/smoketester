# smoketester

![example workflow](https://github.com/smutil/smoketester/actions/workflows/build-actions.yml/badge.svg)![example workflow](https://github.com/smutil/smoketester/actions/workflows/release-actions.yml/badge.svg)

Tool to execute smoke tests. 

Current Features
------------
1. Multiple targets/endpoints can be tested
2. Configurable retry option in case of test failure
3. Configurable Qualitygate
4. Option to provide multiple texts in response for matching
5. GET and POST both method are supported
6. Basic authentication supported


Usage
-----
 step 1. download smoketester from <a href=https://github.com/smutil/smoketester/releases>releases</a>. 
 
 step 2. create [config.yml](https://github.com/smutil/smoketester/blob/main/config.yml). If config.yaml and smoketester is not in same location, you can provide the config.yml path using --config
 
 step 3. execute the smoketester as shown below. 
 
 ```
 ./smoketester --config /path-to-config.yml
 ```

Upcoming Features
------------
1. ssl configuration
2. PostForm method support