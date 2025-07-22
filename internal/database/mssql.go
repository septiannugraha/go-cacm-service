package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/septiannugraha/go-cacm-service/internal/models"
)

// MSSQLClient handles SQL Server database operations
type MSSQLClient struct {
	db *sql.DB
}

// NewMSSQLClient creates a new SQL Server client
func NewMSSQLClient(server, database, user, password string, integratedSecurity bool) (*MSSQLClient, error) {
	var connString string
	if integratedSecurity {
		connString = fmt.Sprintf("server=%s;database=%s;integrated security=true", server, database)
	} else {
		connString = fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", server, database, user, password)
	}

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MSSQLClient{db: db}, nil
}

// Close closes the database connection
func (c *MSSQLClient) Close() error {
	return c.db.Close()
}

// GetSummaryData executes the summary query (similar to GetSummarySiskeudes)
func (c *MSSQLClient) GetSummaryData(tahun string, jenis string) ([]models.CacheData, error) {
	query := c.buildSummaryQuery(tahun, jenis)
	
	rows, err := c.db.Query(query, sql.Named("Tahun", tahun))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []models.CacheData
	for rows.Next() {
		var item models.CacheData
		err := rows.Scan(
			&item.ID,
			&item.Tahun,
			&item.KodeDesa,
			&item.KodeKegiatan,
			&item.KodePaket,
			&item.KodeRekening,
			&item.KodeSumber,
			&item.Tagging,
			&item.Anggaran1,
			&item.Anggaran2,
			&item.Real1,
			&item.Real2,
			&item.Real3,
			&item.Real4,
			&item.Real5,
			&item.Real6,
			&item.Real7,
			&item.Real8,
			&item.Real9,
			&item.Real10,
			&item.Real11,
			&item.Real12,
			&item.TotalReal,
			&item.KodePemda,
			&item.NamaPemda,
			&item.NamaRekening,
			&item.NamaSumber,
			&item.NamaDesa,
			&item.NamaKegiatan,
			&item.NamaPaket,
			&item.IDTipologi,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// buildSummaryQuery builds the SQL query based on year and type
func (c *MSSQLClient) buildSummaryQuery(tahun string, jenis string) string {
	// Check if ID_Tipologi column exists (for 2024+)
	hasIDTipologiColumn := tahun >= "2024"

	var parts []string
	
	// Add base CTEs
	parts = append(parts, SQLQueries.LatestKodePosting)
	parts = append(parts, SQLQueries.KegiatanOutput)
	
	// Add the appropriate CommonSql based on jenis
	switch jenis {
	case "1":
		parts = append(parts, SQLQueries.CommonSQL1) // Revenue
	case "2":
		parts = append(parts, SQLQueries.CommonSQL2) // Expenses
	case "3":
		parts = append(parts, SQLQueries.CommonSQL3) // Financing
	}
	
	// Add tagging CTE
	parts = append(parts, SQLQueries.TaggingSQL)
	
	// Add conditional SQL based on year
	if hasIDTipologiColumn {
		parts = append(parts, SQLQueries.ConditionalSQL2024)
	} else {
		parts = append(parts, SQLQueries.ConditionalSQLBefore2024)
	}
	
	return strings.Join(parts, "\n")
}