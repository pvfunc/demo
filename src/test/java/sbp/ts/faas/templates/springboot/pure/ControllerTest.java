package sbp.ts.faas.templates.springboot.pure;

import io.micrometer.prometheus.PrometheusMeterRegistry;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.web.server.LocalServerPort;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.context.junit4.SpringRunner;

import static org.junit.Assert.assertEquals;

@RunWith(SpringRunner.class)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.DEFINED_PORT)
@TestPropertySource(locations = "classpath:config.yaml")
public class ControllerTest {

    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private PrometheusMeterRegistry registry;

    @Test
    public void testHandle() {
        // given
        String request = "REQUEST";
        String expectedResponse = "Hello from Java11 SpringBoot function!\n" +
                                  "You said: " + request;
        // when
        String response = restTemplate.postForObject("http://localhost:" + port + "/", request, String.class);

        // then
        assertEquals(expectedResponse, response);
        assertEquals(registry.get("custom.template.metric").counter().count(), 1, 0.001);
    }
}