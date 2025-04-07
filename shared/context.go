/*
 * SPDX-License-Identifier: MIT
 * Author: Zenger (https://github.com/Zenger)
 */

package shared

import (
	"TinyBase/config"
	"database/sql"
)

type TinyBaseContext struct {
	Database *sql.DB
	Settings config.Settings
}
