from datadog_checks.base.checks import AgentCheck
import time


class ExampleError(AgentCheck):
    def check(self, instance):
        self.event(
            {
                "timestamp": time.time(),
                "event_type": "Example",
                "msg_title": "Example Error",
                "msg_text": "This is an example error event coming from Datadog.",
                "alert_type": "error",
            }
        )
