#ifndef VOLBLOCK_H
#define VOLBLOCK_H

char *initPulse();
void runPulse();
void setVolume(char, int);
void toggleMute();
void reconnect();

extern void goPulseRestart();
extern void goPulseError(char*);
extern void goPulseMsg(char*);
extern void goPulseVol(char, int);

#endif
