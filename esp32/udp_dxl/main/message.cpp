#include <stdio.h>
#include <stdint.h>
#include "freertos/FreeRTOS.h"
#include "freertos/queue.h"
#include "message.h"


QueueHandle_t message_queue; 

void init_message_queue()
{
	message_queue = xQueueCreate(5, sizeof(msg));
}


