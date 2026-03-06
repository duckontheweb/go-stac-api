#!/bin/bash
pypgstac migrate;
pypgstac load collections --method insert_ignore "/mnt/data/collections.jsonl";
pypgstac load items --method insert_ignore "/mnt/data/naip.jsonl";
pypgstac load items --method insert_ignore "/mnt/data/cop-dem-glo-30.jsonl";
