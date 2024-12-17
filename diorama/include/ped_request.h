#ifndef PED_REQUEST_H
#define PED_REQUEST_H

void button_isr_handler(void *arg); // Remove IRAM_ATTR here
void PedestrianRequestTask(void *pvParameters);
void init_ped_request(void);

#endif // PED_REQUEST_H
