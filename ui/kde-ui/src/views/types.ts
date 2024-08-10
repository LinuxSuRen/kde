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
                image: string;
            };
            redis: {
                enabled: boolean;
                image: string;
            };
            mysql: {
                enabled: boolean;
                image: string;
                password: string;
                database: string;
            };
            postgres: {
                enabled: boolean;
                image: string;
                password: string;
                database: string;
            };
            tDEngine: {
                enabled: boolean;
                image: string;
            };
            rabbitMQ: {
                enabled: boolean;
                image: string;
                username: string;
                password: string;
            };
        };
    };
    status: {
        phase: string;
        deployStatus: string;
        link: string;
    };
}

export function NewEmptyDevSpace() {
    return {
        metadata: {
            name: "",
            namespace: "",
        },
        spec: {
            cpu: "",
            memory: "",
            env: {},
            services: {
                docker: {
                    enabled: false,
                    image: "",
                },
                redis: {
                    enabled: false,
                    image: "",
                },
                mysql: {
                    enabled: false,
                },
                postgres: {
                    enabled: false,
                },
                tDEngine: {
                    enabled: false,
                },
                rabbitMQ: {
                    enabled: false,
                },
            },
        },
    } as DevSpace;
}

export interface Config {
    storageClassName: string;
    volumeMode: string;
    volumeAccessMode: string;
    ingressMode: string;
    imagePullPolicy: string;
    host: string;
}
