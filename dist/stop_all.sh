#!/bin/bash
# Stop all microservices

pkill -f "^\./auth$" || true
pkill -f "^\./catalog$" || true
pkill -f "^\./order$" || true

echo "All services stopped"
