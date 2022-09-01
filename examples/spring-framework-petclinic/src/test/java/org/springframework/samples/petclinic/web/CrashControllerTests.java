/*
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.springframework.samples.petclinic.web;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.junit.jupiter.web.SpringJUnitWebConfig;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.servlet.handler.SimpleMappingExceptionResolver;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/**
 * Test class for {@link CrashController}
 *
 * @author Colin But
 */
@SpringJUnitWebConfig(locations = {"classpath:spring/mvc-core-config.xml", "classpath:spring/mvc-test-config.xml"})
class CrashControllerTests {

    @Autowired
    private CrashController crashController;

    @Autowired
    private SimpleMappingExceptionResolver simpleMappingExceptionResolver;

    private MockMvc mockMvc;

    @BeforeEach
    void setup() {
        this.mockMvc = MockMvcBuilders
            .standaloneSetup(crashController)
            .setHandlerExceptionResolvers(simpleMappingExceptionResolver)
            .build();
    }

    @Test
    void testTriggerException() throws Exception {
        mockMvc.perform(get("/oups"))
            .andExpect(view().name("exception"))
            .andExpect(model().attributeExists("exception"))
            .andExpect(forwardedUrl("exception"))
            .andExpect(status().isOk());
    }
}
