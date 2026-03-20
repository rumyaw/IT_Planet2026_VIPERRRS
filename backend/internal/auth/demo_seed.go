package auth

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureDemoData(ctx context.Context, pool *pgxpool.Pool) error {
	var approvedCount int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM opportunities WHERE status='APPROVED'`).Scan(&approvedCount); err != nil {
		return err
	}
	if approvedCount > 0 {
		return nil
	}

	// Create demo employer.
	demoEmployerEmail := getenvDefault("TRUMPLIN_DEMO_EMPLOYER_EMAIL", "demo_employer@trumplin.local")
	demoEmployerPassword := getenvDefault("TRUMPLIN_DEMO_EMPLOYER_PASSWORD", "demo1234")

	var employerID string
	if err := pool.QueryRow(ctx, `SELECT id::text FROM users WHERE email=$1`, demoEmployerEmail).Scan(&employerID); err != nil {
		pwHash, err := HashPassword(demoEmployerPassword)
		if err != nil {
			return err
		}
		if err := pool.QueryRow(ctx,
			`INSERT INTO users (email, password_hash, role, status, display_name)
			 VALUES ($1,$2,'EMPLOYER','ACTIVE','Demo Employer')
			 RETURNING id::text`,
			demoEmployerEmail, pwHash,
		).Scan(&employerID); err != nil {
			return err
		}
	}

	// Create demo company.
	demoCompanyName := getenvDefault("TRUMPLIN_DEMO_COMPANY_NAME", "IT Planet Demo LLC")
	var companyID string
	if err := pool.QueryRow(ctx,
		`SELECT id::text FROM companies WHERE owner_user_id=$1`,
		employerID,
	).Scan(&companyID); err != nil {
		if err := pool.QueryRow(ctx,
			`INSERT INTO companies (owner_user_id, name, description, verification_status)
			 VALUES ($1,$2,'', 'APPROVED')
			 RETURNING id::text`,
			employerID, demoCompanyName,
		).Scan(&companyID); err != nil {
			return err
		}
	} else {
		_, _ = pool.Exec(ctx, `UPDATE companies SET verification_status='APPROVED' WHERE id=$1`, companyID)
	}

	// Seed tags.
	tagNames := []string{"Go", "PostgreSQL", "React", "SQL", "Docker", "Yandex"}
	tagIDs := make(map[string]string, len(tagNames))

	for _, name := range tagNames {
		// Insert tag if missing.
		_, _ = pool.Exec(ctx, `INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, name)

		var tagID string
		if err := pool.QueryRow(ctx, `SELECT id::text FROM tags WHERE name=$1`, name).Scan(&tagID); err != nil {
			return err
		}
		tagIDs[name] = tagID
	}

	// Seed opportunities.
	city := "Москва"
	opps := []struct {
		title      string
		typ        string
		cityText   string
		lat, lng   float64
		salaryMin  int
		salaryMax  int
		skillNames []string
	}{
		{
			title:      "Go Backend Engineer (Junior+)",
			typ:        "VACANCY",
			cityText:   city,
			lat:         55.7558,
			lng:         37.6173,
			salaryMin:  120000,
			salaryMax:  200000,
			skillNames: []string{"Go", "PostgreSQL", "SQL", "Docker"},
		},
		{
			title:      "Стажировка React + TypeScript",
			typ:        "INTERNSHIP",
			cityText:   city,
			lat:         55.7559,
			lng:         37.6200,
			salaryMin:  80000,
			salaryMax:  120000,
			skillNames: []string{"React", "SQL", "Docker"},
		},
		{
			title:      "Mentor Program: Data + Maps",
			typ:        "MENTOR_PROGRAM",
			cityText:   city,
			lat:         55.7500,
			lng:         37.6050,
			salaryMin:  0,
			salaryMax:  0,
			skillNames: []string{"Go", "Yandex", "PostgreSQL"},
		},
	}

	for _, op := range opps {
		var oppID string
		if err := pool.QueryRow(ctx,
			`INSERT INTO opportunities (
				employer_company_id, curator_user_id,
				title, description, organizer_name,
				type, work_format,
				location_type, city_text, lat, lng,
				salary_min, salary_max,
				status
			) VALUES (
				$1, NULL,
				$2, 'Demo description', $2,
				$3, 'REMOTE',
				'CITY', $4, $5, $6,
				$7, $8,
				'APPROVED'
			) RETURNING id::text`,
			companyID,
			op.title,
			op.typ,
			op.cityText,
			op.lat, op.lng,
			op.salaryMin, op.salaryMax,
		).Scan(&oppID); err != nil {
			return err
		}

		for _, skill := range op.skillNames {
			if tagID := tagIDs[skill]; tagID != "" {
				if _, err := pool.Exec(ctx,
					`INSERT INTO opportunity_tags (opportunity_id, tag_id)
					 VALUES ($1,$2) ON CONFLICT DO NOTHING`,
					oppID, tagID,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

