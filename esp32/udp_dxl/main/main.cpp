#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "esp_system.h"
#include "nvs_flash.h"
#include "esp_netif.h"
#include "esp_event.h"
#include <mdns.h>
#include "lwip/sockets.h"

#include "dxl_task.h"
#include "udp_server_task.h"
#include "protocol_examples_common.h"
#include "message.h"
#include "file_system.h"
#include "esp_spiffs.h"

extern "C" void app_main(void)
{
    ESP_ERROR_CHECK(nvs_flash_init());
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    init_message_queue();

    xTaskCreate(dxl_task, "dynamixel control task", 16384, NULL, 5, NULL);

    initialize_mdns();
}
