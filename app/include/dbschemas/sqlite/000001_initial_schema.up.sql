CREATE TABLE cfg_admin (
    id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    roles TEXT NOT NULL,
    description VARCHAR(256) NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    resource_version VARCHAR(256) NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_cfg_admin_internal_resource_id ON cfg_admin(internal_resource_id);

CREATE TABLE cfg_application (
    id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    chain_id VARCHAR(64) NOT NULL,
    description VARCHAR(256) NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    resource_version VARCHAR(256) NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_cfg_application_internal_resource_id ON cfg_application(internal_resource_id);

CREATE TABLE cfg_user (
    id VARCHAR(64) NOT NULL,
    application_id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    roles TEXT NOT NULL,
    description VARCHAR(256) NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    resource_version VARCHAR(256) NOT NULL,
    PRIMARY KEY (application_id, id)
);
CREATE UNIQUE INDEX idx_cfg_user_internal_resource_id ON cfg_user(internal_resource_id);

CREATE TABLE cfg_account (
    address VARCHAR(64) NOT NULL,
    application_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    PRIMARY KEY (application_id, user_id, address)
);
CREATE UNIQUE INDEX idx_cfg_account_internal_resource_id ON cfg_account(internal_resource_id);

CREATE TABLE cfg_hardware_security_module (
    id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    kind VARCHAR(256) NOT NULL,
    configuration TEXT NOT NULL,
    description VARCHAR(256) NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    resource_version VARCHAR(256) NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_cfg_hardware_security_module_internal_resource_id ON cfg_hardware_security_module(internal_resource_id);

CREATE TABLE cfg_hardware_security_module_slot (
    id VARCHAR(64) NOT NULL,
    internal_resource_id VARCHAR(64) NOT NULL,
    hardware_security_module_id VARCHAR(64) NOT NULL,
    application_id VARCHAR(64) NOT NULL,
    slot VARCHAR(256) NOT NULL,
    pin VARCHAR(256) NOT NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    resource_version VARCHAR(256) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (hardware_security_module_id, slot)
);
CREATE UNIQUE INDEX idx_cfg_hardware_security_module_slot_application_id ON cfg_hardware_security_module_slot(application_id);
CREATE UNIQUE INDEX idx_cfg_hardware_security_module_slot_internal_resource_id ON cfg_hardware_security_module_slot(internal_resource_id);

CREATE TABLE system_referential_integrity_entry (
    id VARCHAR(64) NOT NULL,
    resource_id VARCHAR(64) NOT NULL,
    resource_kind VARCHAR(256) NOT NULL,
    parent_resource_id VARCHAR(64) NOT NULL,
    parent_resource_kind VARCHAR(256) NOT NULL,
    creation_date BIGINT NULL,
    last_update BIGINT NULL,
    PRIMARY KEY (id),
    UNIQUE (resource_id, resource_kind, parent_resource_id, parent_resource_kind)
);
CREATE INDEX idx_system_referential_integrity_entry_resource ON system_referential_integrity_entry (resource_id, resource_kind);
CREATE INDEX idx_system_referential_integrity_entry_referenced_resource ON system_referential_integrity_entry (parent_resource_id, parent_resource_kind);
