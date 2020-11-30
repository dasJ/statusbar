#include <pulse/pulseaudio.h>
#include <string.h>

#include "volblock.h"

pa_context *context;
pa_mainloop *mainloop = NULL;
int oldMute;
int connected = 0;
pa_cvolume oldVolume;
char *defaultSinkName;

void sink_info_cb(pa_context *c, const pa_sink_info *i, int eol, void *userdata) {
	if (i == NULL || defaultSinkName == NULL) {
		return;
	}
	if (strcmp(i->name, defaultSinkName)) {
		return;
	}

	oldVolume = i->volume;
	oldMute = i->mute;
	goPulseVol(i->mute, ((float) pa_cvolume_avg(&i->volume) / PA_VOLUME_NORM) * 100);
}

void server_info_cb(pa_context *c, const pa_server_info *i, void *userdata) {
	pa_operation *o;

	if (defaultSinkName == NULL) {
		free(defaultSinkName);
	}
	// Update name of the default sink
	if (!i->default_sink_name) {
		return;
	}
	defaultSinkName = malloc(sizeof(char) * strlen(i->default_sink_name) + 1);
	strcpy(defaultSinkName, i->default_sink_name);
	// Request info of the (new) default sink
	if (!(o = pa_context_get_sink_info_by_name(c, defaultSinkName, sink_info_cb, NULL))) {
		goPulseError("Cannot request sink info");
		return;
	}
}

void subscription_event_cb(pa_context *c, pa_subscription_event_type_t t, uint32_t idx, void *userdata) {
	pa_operation *o;

	switch (t & PA_SUBSCRIPTION_EVENT_FACILITY_MASK) {
		case PA_SUBSCRIPTION_EVENT_SINK:
			if (!(o = pa_context_get_sink_info_by_index(c, idx, sink_info_cb, NULL))) {
				goPulseError("Cannot get sink info");
				return;
			}
			pa_operation_unref(o);
			break;
		case PA_SUBSCRIPTION_EVENT_SERVER:
			if (!(o = pa_context_get_server_info(c, server_info_cb, NULL))) {
				goPulseError("Cannot get server info");
				return;
			}
			break;
		default:
			break;
	}
}

void context_state_cb(pa_context *c, void *userdata) {
	pa_operation *o;
	int ret;

	switch (pa_context_get_state(c)) {
		case PA_CONTEXT_READY:
			// Register our callback
			pa_context_set_subscribe_callback(c, subscription_event_cb, NULL);
			// Request server info
			if (!(o = pa_context_get_server_info(c, server_info_cb, NULL))) {
				goPulseError("Cannot request PulseAudio server info");
				return;
			}
			pa_operation_unref(o);
			// Subscribe to events
			if (!(o = pa_context_subscribe(c, (pa_subscription_mask_t) (PA_SUBSCRIPTION_MASK_SINK | PA_SUBSCRIPTION_MASK_SERVER), NULL, NULL))) {
				goPulseError("Cannot subscribe to PulseAudio events");
				return;
			}
			pa_operation_unref(o);
			connected = 1;
			break;
		case PA_CONTEXT_FAILED:
			goPulseMsg("Disconnected");
			connected = 0;
			reconnect();
		default:
			break;
	}
}

char *initPulse() {
	pa_mainloop_api *mainloop_api;
	int ret = 0;

	// Initialize main loop
	if (!(mainloop = pa_mainloop_new())) {
		return "Cannot create mainloop";
	}
	mainloop_api = pa_mainloop_get_api(mainloop);

	// Initialize context
	if (!(context = pa_context_new(mainloop_api, "i3-statusbar"))) {
		return "Cannot set up context";
	}

	// Context state callback
	pa_context_set_state_callback(context, context_state_cb, NULL);

	// Connect the context
	if (pa_context_connect(context, NULL, 0, NULL) < 0) {
		return "Cannot connect context";
	}

	return "";
}

void reconnect() {
	char *err;

	if (context)
		pa_context_unref(context);
	if (mainloop)
		pa_mainloop_free(mainloop);

	goPulseRestart();
}

void runPulse() {
	int ret = 0;

	// Run the main loop
	if (pa_mainloop_run(mainloop, &ret) < 0) {
		goPulseError("Cannot start main loop");
	}
}

void setVolume(char increase, int amount) {
	pa_cvolume *newVolume;

	if (!connected)
		return;

	if (increase == 1) {
		newVolume = pa_cvolume_inc(&oldVolume, amount * ((float) PA_VOLUME_NORM / 100));
	} else {
		newVolume = pa_cvolume_dec(&oldVolume, amount * ((float) PA_VOLUME_NORM / 100));
	}

	pa_operation_unref(pa_context_set_sink_volume_by_name(context, defaultSinkName, newVolume, NULL, NULL));
}

void toggleMute() {
	if (!connected)
		return;

	pa_operation_unref(pa_context_set_sink_mute_by_name(context, defaultSinkName, !oldMute, NULL, NULL));
}
