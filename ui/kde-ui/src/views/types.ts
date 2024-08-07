/*
Copyright 2024 kde authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

export interface DevSpace {
    metadata: {
        name: string;
        namespace: string;
    };
    spec: {
        cpu: string;
        memory: string;
        env: {};
        services: {
            docker: {
                enabled: boolean;
            };
            redis: {
                enabled: boolean;
            };
            mysql: {
                enabled: boolean;
            };
        };
    };
    status: {
        phase: string;
        deployStatus: string;
        link: string;
    };
}
