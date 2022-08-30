/*
 * Copyright 2002-2013 the original author or authors.
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
package org.springframework.samples.petclinic.repository.jdbc;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.namedparam.MapSqlParameterSource;
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcTemplate;
import org.springframework.jdbc.core.simple.SimpleJdbcInsert;
import org.springframework.samples.petclinic.model.Visit;
import org.springframework.samples.petclinic.repository.VisitRepository;
import org.springframework.stereotype.Repository;

import javax.sql.DataSource;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * A simple JDBC-based implementation of the {@link VisitRepository} interface.
 *
 * @author Ken Krebs
 * @author Juergen Hoeller
 * @author Rob Harrop
 * @author Sam Brannen
 * @author Thomas Risberg
 * @author Mark Fisher
 * @author Michael Isvy
 */
@Repository
public class JdbcVisitRepositoryImpl implements VisitRepository {

    private NamedParameterJdbcTemplate jdbcTemplate;

    private SimpleJdbcInsert insertVisit;

    @Autowired
    public JdbcVisitRepositoryImpl(DataSource dataSource) {
        this.jdbcTemplate = new NamedParameterJdbcTemplate(dataSource);

        this.insertVisit = new SimpleJdbcInsert(dataSource)
            .withTableName("visits")
            .usingGeneratedKeyColumns("id");
    }


    @Override
    public void save(Visit visit) {
        if (visit.isNew()) {
            Number newKey = this.insertVisit.executeAndReturnKey(
                createVisitParameterSource(visit));
            visit.setId(newKey.intValue());
        } else {
            throw new UnsupportedOperationException("Visit update not supported");
        }
    }


    /**
     * Creates a {@link MapSqlParameterSource} based on data values from the supplied {@link Visit} instance.
     */
    private MapSqlParameterSource createVisitParameterSource(Visit visit) {
        return new MapSqlParameterSource()
            .addValue("id", visit.getId())
            .addValue("visit_date", visit.getDate())
            .addValue("description", visit.getDescription())
            .addValue("pet_id", visit.getPet().getId());
    }

    @Override
    public List<Visit> findByPetId(Integer petId) {
        Map<String, Object> params = new HashMap<>();
        params.put("id", petId);
        JdbcPet pet = this.jdbcTemplate.queryForObject(
                "SELECT id, name, birth_date, type_id, owner_id FROM pets WHERE id=:id",
                params,
                new JdbcPetRowMapper());

        List<Visit> visits = this.jdbcTemplate.query(
            "SELECT id as visit_id, visit_date, description FROM visits WHERE pet_id=:id",
            params, new JdbcVisitRowMapper());

        for (Visit visit: visits) {
            visit.setPet(pet);
        }

        return visits;
    }

}
