package sbp.ts.faas.templates.springboot.pure;

import io.micrometer.core.instrument.Tags;
import io.micrometer.prometheus.PrometheusMeterRegistry;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class Controller {
    private static final Logger logger = LoggerFactory.getLogger(Controller.class);

    private final PrometheusMeterRegistry registry;

    @Autowired
    public Controller(PrometheusMeterRegistry registry) {
        this.registry = registry;
    }

    @PostMapping(value = "/")
    public ResponseEntity<String> handle(@RequestBody String payload) {
        // Log request
        logger.info("Function called with request: {}", payload);
        // Send custom metric
        registry.counter("custom.template.metric", Tags.empty()).increment();

        return ResponseEntity.ok()
            .contentType(MediaType.TEXT_PLAIN)
            .body(
                "Hello from Java11 SpringBoot function!\n" +
                "You said: " + payload
            );
    }
}
