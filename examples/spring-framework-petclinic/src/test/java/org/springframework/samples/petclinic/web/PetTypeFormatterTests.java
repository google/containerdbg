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

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.samples.petclinic.model.PetType;
import org.springframework.samples.petclinic.service.ClinicService;

import java.text.ParseException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Locale;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * Test class for {@link PetTypeFormatter}
 *
 * @author Colin But
 */
@ExtendWith(MockitoExtension.class)
class PetTypeFormatterTests {

    @Mock
    private ClinicService clinicService;

    private PetTypeFormatter petTypeFormatter;

    @BeforeEach
    void setup() {
        petTypeFormatter = new PetTypeFormatter(clinicService);
    }

    @Test
    void testPrint() {
        PetType petType = new PetType();
        petType.setName("Hamster");
        String petTypeName = petTypeFormatter.print(petType, Locale.ENGLISH);
        assertEquals("Hamster", petTypeName);
    }

    @Test
    void shouldParse() throws ParseException {
        Mockito.when(clinicService.findPetTypes()).thenReturn(makePetTypes());
        PetType petType = petTypeFormatter.parse("Bird", Locale.ENGLISH);
        assertEquals("Bird", petType.getName());
    }

    @Test
    void shouldThrowParseException() throws ParseException {
        Mockito.when(clinicService.findPetTypes()).thenReturn(makePetTypes());
        Assertions.assertThrows(ParseException.class, () -> {
            petTypeFormatter.parse("Fish", Locale.ENGLISH);
        });
    }

    /**
     * Helper method to produce some sample pet types just for test purpose
     *
     * @return {@link Collection} of {@link PetType}
     */
    private Collection<PetType> makePetTypes() {
        Collection<PetType> petTypes = new ArrayList<>();
        petTypes.add(new PetType(){
            {
                setName("Dog");
            }
        });
        petTypes.add(new PetType(){
            {
                setName("Bird");
            }
        });
        return petTypes;
    }

}
