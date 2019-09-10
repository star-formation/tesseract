/*  Copyright 2019 The tesseract Authors

    This file is part of tesseract.

    tesseract is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    tesseract is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package tesseract

import (
	"github.com/ethereum/go-ethereum/log"

	"database/sql"

	_ "github.com/lib/pq"

    "golang.org/x/image/math/fixed"
)

var db *sql.DB

var dropSchema = `
DROP SCHEMA IF EXISTS an CASCADE
`

var createSchema = `
CREATE SCHEMA an
    CREATE TABLE entity (id bigserial UNIQUE)

	CREATE TABLE spatial (entity bigint,
        p_x bigint, p_y bigint, p_z bigint,
        v_x bigint, v_y bigint, v_z bigint, v_m bigint,
        o_x bigint, o_y bigint, o_z bigint,
        r_x bigint, r_y bigint, r_z bigint, r_m bigint)
	
	CREATE TABLE physical (entity bigint, mass bigint)
	
	CREATE TABLE orbit (entity bigint)
	
	CREATE TABLE spherical_shape (entity bigint, radius bigint)
`

type DBSpatial struct {
    // Entity id
    entity int64
    // 3D position
    p_x, p_y, p_z fixed.Int52_12
    // 3D velocity: direction of speed and its magnitude
    v_x, v_y, v_z, v_m fixed.Int52_12 
    // 3D orientation: direction of object vs reference frame
    o_x, o_y, o_z fixed.Int52_12
    // 3D rotation: angular velocity; direction of rotation and its magnitude
    r_x, r_y, r_z, r_m fixed.Int52_12
}

type DBPhysical struct {
    mass fixed.Int52_12
}

var newEntityStmt *sql.Stmt
var querySpatialStmt *sql.Stmt
var insertSpatialStmt *sql.Stmt
var insertSphericalShapeStmt *sql.Stmt

func init() {
    connStr := "user=argo_navis_test password=test dbname=argo_navis_test sslmode=disable"
    d, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Error("sql.Open", err, err)
    }
    db = d

    // TODO: refactor drop / create, this is for testing
    _, err = query(dropSchema)
    if err != nil { panic(err) }
    _, err = query(createSchema)
    if err != nil { panic(err) }

    prepareStmt := func(s **sql.Stmt, q string) {
        stmt, err := db.Prepare(q)
        if err != nil {
            log.Error("db.Prepare", "err", err)
            panic(err)
        }
        *s = stmt
    }

    prepareStmt(&querySpatialStmt, "SELECT entity, p_x, p_y, p_z, v_x, v_y, v_z, v_m, o_x, o_y, o_z, r_x, r_y, r_z, r_m FROM an.spatial")
    prepareStmt(&insertSpatialStmt, "INSERT INTO an.spatial (entity, p_x, p_y, p_z, v_x, v_y, v_z, v_m, o_x, o_y, o_z, r_x, r_y, r_z, r_m) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)")
    prepareStmt(&newEntityStmt, "INSERT INTO an.entity (id) VALUES(DEFAULT) RETURNING id")
    prepareStmt(&insertSphericalShapeStmt, "INSERT INTO an.spherical_shape (entity, radius) VALUES($1, $2)")
}

func querySpatial() ([]*DBSpatial) {
    rows, err := querySpatialStmt.Query()
    if err != nil {
        log.Error("querySpatialStmt.Query", "err", err)
        return nil
    }
    defer rows.Close() // TODO: refactor overall DB lifecycle / error handling

    res := make([]*DBSpatial, 1)
    for rows.Next() {
        dest := new(DBSpatial)
        err := rows.Scan(&dest.entity, &dest.p_x, &dest.p_y, &dest.p_z, &dest.v_x, &dest.v_y, &dest.v_z, &dest.v_m, &dest.o_x, &dest.o_y, &dest.o_z, &dest.r_x, &dest.r_y, &dest.r_z, &dest.r_m)
        if err != nil{
            log.Error("sql.Rows.Scan", "err", err)
            return nil
        }
        res = append(res, dest)
    }
    
    err = rows.Err()
    if err != nil{
        log.Error("sql.Rows.Err", "err", err)
        return nil
    }

    return res
}

func insertSpatial(s *DBSpatial) error {
    _, err := insertSpatialStmt.Exec(s.entity,
                                     s.p_x, s.p_y, s.p_z,
                                     s.v_x, s.v_y, s.v_z, s.v_m,
                                     s.o_x, s.o_y, s.o_z,
                                     s.r_x, s.r_y, s.r_z, s.r_m)
    if err != nil {
        log.Error("insertSpatialStmt.Exec", "err", err)
        return err
    }
    return nil
}

func newEntity() (int64, error) {
    var id int64
    err := newEntityStmt.QueryRow().Scan(&id)
    if err != nil {
        log.Error("newEntityStmt.QueryRow", "err", err)
        return 0, err
    }
    return id, nil
}

func insertSphericalShape(entity, radius int64) error {
    _, err := insertSphericalShapeStmt.Exec(entity, radius)
    if err != nil {
        log.Error("insertSphericalShapeStmt", "err", err)
        return err
    }
    return nil
}

func query(s string) (*sql.Rows, error) {
    rows, err := db.Query(s)
    if err != nil{
        log.Error("sql.DB.Query", "err", err)
        return nil, err
    }
    return rows, err
}
