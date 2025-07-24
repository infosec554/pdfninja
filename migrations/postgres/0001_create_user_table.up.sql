CREATE TYPE user_status_enum AS ENUM ('active', 'deleted', 'blocked');
CREATE TYPE otp_status_enum AS ENUM ('unconfirmed', 'confirmed');

-- 2. Table: users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    status user_status_enum NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID
);

-- 3. Table: otp
CREATE TABLE otp (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    status otp_status_enum NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP
);

-- 4. Table: sysusers
CREATE TABLE sysusers (
    id UUID PRIMARY KEY,
    status user_status_enum NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID
);

-- 5. Table: roles
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    status user_status_enum NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID
);

-- 6. Table: sysuser_roles (many-to-many)
CREATE TABLE sysuser_roles (
    id UUID PRIMARY KEY,
    sysuser_id UUID NOT NULL REFERENCES sysusers(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
);


--************************************************ FILE STORAGE
CREATE TABLE files (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(20) NOT NULL,
    file_size INTEGER NOT NULL,
    uploaded_at TIMESTAMP DEFAULT NOW()
);

-- ORGANIZE PDF

CREATE TABLE organize_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    new_order INTEGER[] NOT NULL, -- Sahifalar yangi tartibi
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE merge_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE merge_job_input_files (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES merge_jobs(id) ON DELETE CASCADE,
    file_id UUID NOT NULL REFERENCES files(id)
);

CREATE TABLE split_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    split_ranges TEXT NOT NULL,
    output_file_ids UUID[] DEFAULT ARRAY[]::UUID[],
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE remove_pages_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    pages_to_remove TEXT NOT NULL,
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE extract_pages_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    pages_to_extract TEXT NOT NULL,
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- OPTIMIZE PDF

CREATE TABLE compress_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    compression_level VARCHAR(10),
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- CONVERT TO PDF

CREATE TABLE jpg_to_pdf_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_ids UUID[] NOT NULL, -- Bir nechta rasm fayllari
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE word_to_pdf_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id), -- Bitta Word fayl
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE pdf_to_jpg_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    output_file_ids UUID[] DEFAULT ARRAY[]::UUID[],
    zip_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);



-- EDIT PDF

CREATE TABLE rotate_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    rotation_angle INTEGER NOT NULL,
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE rotate_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    rotation_angle INTEGER NOT NULL,
    pages VARCHAR NOT NULL DEFAULT 'all',
    output_file_id UUID REFERENCES files(id),
    output_path VARCHAR,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE add_page_number_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    style VARCHAR(50),
    position VARCHAR(20),
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE watermark_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    watermark_text TEXT NOT NULL,
    position VARCHAR(20),
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE crop_pdf_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    top INTEGER NOT NULL,
    bottom INTEGER NOT NULL,
    "left" INTEGER NOT NULL,
    "right" INTEGER NOT NULL,
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- PDF SECURITY

CREATE TABLE unlock_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE protect_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    password TEXT NOT NULL,
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE logs (
    id UUID PRIMARY KEY,
    job_id UUID,  -- merge_jobs.id, split_jobs.id va h.k.
    job_type VARCHAR(30), -- masalan: 'merge', 'split', 'compress'
    message TEXT,
    level VARCHAR(10), -- 'info', 'error', 'debug'
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE pdf_inspect_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    file_id UUID NOT NULL REFERENCES files(id),
    page_count INT,
    title TEXT,
    author TEXT,
    subject TEXT,
    keywords TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE translate_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    source_lang VARCHAR(10) NOT NULL,         -- masalan: 'en'
    target_lang VARCHAR(10) NOT NULL,         -- masalan: 'uz'
    output_file_id UUID REFERENCES files(id), -- tarjima qilingan PDF fayl
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE shared_links (
    id UUID PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    shared_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE add_header_footer_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    
    header_text TEXT,
    footer_text TEXT,
    font_size INTEGER DEFAULT 12,
    font_color VARCHAR(20) DEFAULT 'black',
    position VARCHAR(20) DEFAULT 'center', -- left, center, right

    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE add_background_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    background_image_file_id UUID NOT NULL REFERENCES files(id),
    opacity FLOAT DEFAULT 1.0,  -- 0.0 dan 1.0 gacha
    position VARCHAR(50) DEFAULT 'center', -- masalan: 'center', 'top-left', 'bottom-right'
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE detect_blank_pages_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    blank_pages INTEGER[] DEFAULT ARRAY[]::INTEGER[], -- bo‘sh sahifa raqamlari saqlanadi
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE qr_code_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    qr_content TEXT NOT NULL, -- QR kod uchun matn
    position VARCHAR(20) NOT NULL, -- qr kod joylashuvi (masalan: "top-left", "center")
    size INTEGER NOT NULL, -- QR kod o'lchami pikselda
    output_file_id UUID REFERENCES files(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE pdf_text_search_jobs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    input_file_id UUID NOT NULL REFERENCES files(id),
    extracted_text TEXT, -- Fayldan chiqarilgan matn
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);
