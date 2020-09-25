

-- ----------------------------
-- Table structure for master
-- ----------------------------
DROP TABLE IF EXISTS "public"."master";
CREATE TABLE "public"."master" (
"id" varchar(100) COLLATE "default" NOT NULL,
"name" varchar(100) COLLATE "default" NOT NULL,
"password" varchar(100) COLLATE "default",
"mobile" varchar(100) COLLATE "default"
)
WITH (OIDS=FALSE)

;
COMMENT ON TABLE "public"."master" IS '帐号表';

-- ----------------------------
-- Alter Sequences Owned By 
-- ----------------------------
