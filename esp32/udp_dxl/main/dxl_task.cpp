#include "dxl_task.h"
#include "dynamixel2espressif.h"
#include "esp_log.h"
#include "message.h"
#include "file_system.h"

const float DXL_PROTOCOL_VERSION = 2.0;
static const char * TAG = "DXL";

#define NUM_JOINTS 18

const int COXA_INDEX = 1;
const int FEMUR_INDEX = 2;
const int TIBIA_INDEX = 3;

Dynamixel2Espressif dxl(UART_NUM_2, (gpio_num_t)CONFIG_DXL_DIR_PIN);

void setup_dxl() {
  // put your setup code here, to run once:
  
  // Set Port baudrate to 1000000bps. This has to match with DYNAMIXEL baudrate.
  dxl.begin(1000000);
  // Set Port Protocol Version. This has to match with DYNAMIXEL protocol version.
  dxl.setPortProtocolVersion(DXL_PROTOCOL_VERSION);

  // Turn off torque when configuring items in EEPROM area

  // We're assuming that servo IDs are from 1 to 18
  // All servos are configured to joint / position mode
  for (int servoId=1; servoId<=NUM_JOINTS; servoId++)
  {

    ESP_LOGI(TAG, "Configuring servo %d", servoId);

    while (0 == dxl.torqueOff(servoId));
    while (0 == dxl.setOperatingMode(servoId, OP_POSITION));
    while (0 == dxl.torqueOn(servoId));

    // if (0 == dxl.ledOn(servoId))
    // {
    //   ESP_LOGE(TAG, "Failed to switch LED on for id: %d", servoId);
    // }
  }
}


void process_stream() 
{
    msg udp_message;
    extern QueueHandle_t message_queue; 

    while (true )
    {
      if(xQueueReceive(message_queue, &(udp_message) , (TickType_t)5 ))
      { 
        // ESP_LOGI(TAG, "------------------------");
        for (int i=0; i<NUM_LEGS; i++)
        {
          dxl.setGoalPosition(i*3+COXA_INDEX, udp_message.packet.legs[i].coxa);
          dxl.setGoalPosition(i*3+FEMUR_INDEX, udp_message.packet.legs[i].femur);
          dxl.setGoalPosition(i*3+TIBIA_INDEX, udp_message.packet.legs[i].tibia);

          // ESP_LOGI(TAG, "Leg: %d, Coxa: %d, Femur: %d, Tibia: %d", i+1, udp_message.packet.legs[i].coxa, udp_message.packet.legs[i].femur, udp_message.packet.legs[i].tibia);
        }
      }
    }
    vTaskDelay(1);
}

int buffer16[5000];
uint8_t buffer8[5000];


char path[64];

void run_primitive(char * name)
{
  sprintf(path,"/spiffs/%s", name);
  FILE* f = fopen(path, "rb");
  if (f == NULL) {
    ESP_LOGE(TAG, "Failed to open file for reading");
    return;
  }

  int nread = fread(buffer8, 1, sizeof(buffer8), f);

  int n16 = 0;
   for (int i=0; i<nread; i+=2) {
      buffer16[n16] = (buffer8[i] & 0xFF) | (buffer8[i+1] << 8);
      n16++;
   }

    for (int step =0; step<nread/2; step += 18) {
      for (int servo=0; servo <18; servo++) {
        while(!dxl.setGoalPosition(servo+1, buffer16[step+servo]));
      }
       vTaskDelay(2);
    }
}

void dxl_task(void *pvParameters)
{
  setup_dxl();

  initialize_spiffs();

  // while (true) 
  // {
    run_primitive("tripod");
    run_primitive("rotate");


    // run_primitive("rotate");
  // }

  while (true) 
  {
    vTaskDelay(100);

  }
}