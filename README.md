# smoketester

![example workflow](https://github.com/smutil/smoketester/actions/workflows/build-actions.yml/badge.svg)![example workflow](https://github.com/smutil/smoketester/actions/workflows/release-actions.yml/badge.svg)

Tool to execute smoke tests. 

Current Features
------------
1. Multiple targets/enpoints can be tested
2. Configurable retry option in case of test failure
3. Configurable Qualitygate


Usage
-----
 step 1. download smoketester from <a href=https://github.com/smutil/smoketester/releases>releases</a>. 
 
 step 2. create [config.yml](https://github.com/smutil/smoketester/config.yml). If config.yaml and smoketester is not in same location, you can provide the config.yml path using --config
 
 step 3. execute the smoketester as shown below. 
 
 ```
 ./smoketester --config /path-to-config.yml
 ```

Upcoming Features
------------
1. ssl configuration
2. Post method support