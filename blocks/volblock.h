#ifndef VOLBLOCK_H
#define VOLBLOCK_H

char *initPulse();
void runPulse();
void setVolume(char, int);
void toggleMute();

extern void goPulseError(char*);
extern void goPulseVol(char, int);

#endif
