#!/bin/bash

sqlboiler mysql --config config/common.toml --output internal/db/model --pkgname model
