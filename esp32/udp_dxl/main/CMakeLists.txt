idf_component_register(
                    SRCS 
                        "main.cpp"
                        "udp_server_task.cpp"
                        "dxl_task.cpp"
                        "message.cpp"
                        "file_system.cpp"
                    INCLUDE_DIRS 
                        "." 
                    REQUIRES 
                        esp_wifi
                        nvs_flash
                        dynamixel2espressif
                        spiffs
                    )

                    spiffs_create_partition_image(spiffs ../../../goik/primitives FLASH_IN_PROJECT)