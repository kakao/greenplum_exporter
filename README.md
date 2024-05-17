# Greenplum Exporter

This is for Greenplum 6 and 7

**Pre-requisites**

go1.21

### 1. Build a greenplum exporter

1.1. git clone

1.2. build an exporter

*Before run it, you need to edit the Makefile for your environment*

```
make build
```


### 2. Run a greenplum exporter
Set the environment variable, GPDB_CONNECTION.

**Connection String**
- The GPDB_CONNECTION is as following form:
```
postgres://[ID]:[PASSWORD]@[ADDRESS]:[PORT]/[DATABASE]?[OPTION]=[VALUE]&[OPTION]=[VALUE]...
```

Should consider the **sslmode** parameter

*Example*
```
export GPDB_CONNECTION=postgres://gpadmin:password@localhost:5432/postgres?sslmode=disable
greenplum_exporter --loglevel=error --greenplumVersion=7 --excludeScraps="size_detail"
```

The exporter displays metrics at *http://IP:9101/metrics*

**Options**
```
--web.listen-address="0.0.0.0:9101"
                        web endpoint
--excludeScraps=""      excluding metric names: user,segment,conn,lock,size,size_detail,bgwriter
--greenplumVersion="6"  greenplum Version, greenplum 6.x, greenplum 7.x: 6, 7
--loglevel="info"       Set the level of log: debug, info, warn, error
--[no-]version          Show application version.
```

## License

This software is licensed under the [Apache 2 license](LICENSE), quoted below.

Copyright 2024 Kakao Corp. <http://www.kakaocorp.com>

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this project except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
