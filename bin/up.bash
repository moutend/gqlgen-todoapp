#!/bin/bash

migrate -database "mysql://root:abcdef123456@(127.0.0.1:33306)/todoapp" -path migrations/mysql up
