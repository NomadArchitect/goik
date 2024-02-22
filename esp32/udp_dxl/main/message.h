#ifndef _MESSAGE_H_
#define _MESSAGE_H_

#define MESSAGE_LENGTH 39
#define NUM_LEGS 6

#pragma pack(push, 1)
typedef union
{
	uint8_t rx_buffer[MESSAGE_LENGTH];
	struct 
	{
		uint8_t networkId;
		struct 
		{
			uint16_t coxa;
			uint16_t femur;
			uint16_t tibia;
		}	legs[NUM_LEGS];
	} packet;
} msg;
#pragma pack(pop)

void init_message_queue();


#endif // _MESSAGE_H_