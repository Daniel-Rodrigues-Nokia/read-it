package com.example.demo;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Nested;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.hamcrest.Matchers.containsString;
import static org.hamcrest.Matchers.equalTo;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
public class ExampleTests {

    @Autowired
    private MockMvc mockMvc;

    @BeforeEach
    void setup() {
        // runs before each test
    }

    @Test
    void helloEndpointShouldReturnOk() throws Exception {
        mockMvc.perform(get("/hello"))
            .andExpect(status().isOk())
            .andExpect(content().string("Hello World"));
    }

    @Test
    void shouldCreateResource() throws Exception {
        String json = "{\"name\":\"John\",\"age\":30}";

        mockMvc.perform(post("/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(json))
            .andExpect(status().isCreated());
    }

    @Test
    void shouldMockServiceCall() throws Exception {
        mockMvc.perform(get("/users/1"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.name", equalTo("Mock User")));
    }

    @Test
    void shouldReturnListOfUsers() throws Exception {
        mockMvc.perform(get("/users"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.length()").value(3));
    }

    @Test
    void shouldReturnCustomHeader() throws Exception {
        mockMvc.perform(get("/hello"))
            .andExpect(header().exists("X-Custom-Header"));
    }

    // ✅ 6. Data-driven test (like Cypress .forEach)
    @Nested
    class MultipleInputTests {
        @Test
        void testMultipleInputs() throws Exception {
            String[][] cases = {
                {"Test simple string", "hello"},
                {"Test numeric string", "12345"},
                {"Test random string", "JUnitRocks"}
            };

            for (String[] testCase : cases) {
                String desc = testCase[0];
                String value = testCase[1];

                mockMvc.perform(post("/echo")
                        .contentType(MediaType.TEXT_PLAIN)
                        .content(value))
                    .andExpect(status().isOk())
                    .andExpect(content().string(containsString(value)));
            }
        }
    }
}
