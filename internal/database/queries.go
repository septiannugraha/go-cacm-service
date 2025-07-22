package database

// SQLQueries contains all SQL query constants
var SQLQueries = struct {
	LatestKodePosting         string
	KegiatanOutput           string
	CommonSQL1               string
	CommonSQL2               string
	CommonSQL3               string
	TaggingSQL               string
	ConditionalSQL2024       string
	ConditionalSQLBefore2024 string
}{
	LatestKodePosting: `
    SET ARITHABORT ON;
        WITH LatestKdPosting AS (
            SELECT MAX(Q.KdPosting) AS MaxKdPosting, Q.Tahun, Q.Kd_Desa
            FROM Ta_AnggaranLog AS Q
            GROUP BY Q.Tahun, Q.Kd_Desa
        ),
    `,

	KegiatanOutput: `
        KegiatanOutput AS (
            SELECT TOP 1 No_ID, Kd_Desa, Kd_Keg, Nama_Paket 
            FROM Ta_KegiatanOutput 
            WHERE AnggaranPAK IS NOT NULL
        ),
    `,

	CommonSQL1: `
				CTE_AggregatedData AS (
    SELECT
        Tahun,
        Kd_Desa,
        Kd_Keg,
        COALESCE(Kd_SubRinci, '01') AS Kd_SubRinci,
        Kd_Rincian,
        Sumberdana,
		SUM(CASE WHEN KdPosting = 2 THEN Anggaran ELSE 0 END) AS Anggaran1,
        SUM(CASE WHEN KdPosting = MaxKdPosting THEN (CASE WHEN KdPosting = 2 THEN Anggaran ELSE AnggaranStlhPAK END) ELSE 0 END) AS Anggaran2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 1 THEN Nilai ELSE 0 END) AS Real1,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 2 THEN Nilai ELSE 0 END) AS Real2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 3 THEN Nilai ELSE 0 END) AS Real3,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 4 THEN Nilai ELSE 0 END) AS Real4,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 5 THEN Nilai ELSE 0 END) AS Real5,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 6 THEN Nilai ELSE 0 END) AS Real6,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 7 THEN Nilai ELSE 0 END) AS Real7,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 8 THEN Nilai ELSE 0 END) AS Real8,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 9 THEN Nilai ELSE 0 END) AS Real9,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 10 THEN Nilai ELSE 0 END) AS Real10,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 11 THEN Nilai ELSE 0 END) AS Real11,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 12 THEN Nilai ELSE 0 END) AS Real12,
        SUM(Nilai) AS TotalReal
    FROM (
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, Anggaran, AnggaranStlhPAK, 0 AS Nilai, NULL AS Tgl_Bukti, A.KdPosting, LKP.MaxKdPosting FROM Ta_AnggaranRinci A
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
       LEFT JOIN LatestKdPosting LKP ON A.Tahun = LKP.Tahun AND A.Kd_Desa = LKP.Kd_Desa
        WHERE LEFT(A.Kd_Rincian, 1) = '4'
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_TBPRinci A
        INNER JOIN Ta_TBP ON A.No_Bukti = Ta_TBP.No_Bukti
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 1) = '4'
        AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_Mutasi A
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 1) = '4'
          AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (A.Kredit - A.Debet) Nilai, NULL AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_JurnalUmumRinci A
        INNER JOIN Ta_JurnalUmum B ON A.NoBukti = B.NoBukti
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 1) = '4' AND B.Posted = 1
    ) AS CombinedTables
    GROUP BY Tahun, Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana
),
`,

	CommonSQL2: `
				CTE_AggregatedData AS (
    SELECT
        Tahun,
        Kd_Desa,
        Kd_Keg,
        COALESCE(Kd_SubRinci, '01') AS Kd_SubRinci,
        Kd_Rincian,
        Sumberdana,
		        SUM(CASE WHEN KdPosting = 2 THEN Anggaran ELSE 0 END) AS Anggaran1,
        SUM(CASE WHEN KdPosting = MaxKdPosting THEN AnggaranStlhPAK ELSE 0 END) AS Anggaran2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 1 THEN Nilai ELSE 0 END) AS Real1,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 2 THEN Nilai ELSE 0 END) AS Real2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 3 THEN Nilai ELSE 0 END) AS Real3,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 4 THEN Nilai ELSE 0 END) AS Real4,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 5 THEN Nilai ELSE 0 END) AS Real5,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 6 THEN Nilai ELSE 0 END) AS Real6,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 7 THEN Nilai ELSE 0 END) AS Real7,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 8 THEN Nilai ELSE 0 END) AS Real8,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 9 THEN Nilai ELSE 0 END) AS Real9,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 10 THEN Nilai ELSE 0 END) AS Real10,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 11 THEN Nilai ELSE 0 END) AS Real11,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 12 THEN Nilai ELSE 0 END) AS Real12,
        SUM(Nilai) AS TotalReal
    FROM (
                SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, Anggaran, AnggaranStlhPAK, 0 AS Nilai, NULL AS Tgl_Bukti, A.KdPosting, LKP.MaxKdPosting FROM Ta_AnggaranRinci A
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
       LEFT JOIN LatestKdPosting LKP ON A.Tahun = LKP.Tahun AND A.Kd_Desa = LKP.Kd_Desa
        WHERE LEFT(A.Kd_Rincian, 1) = '5'
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, C.Tgl_Cek AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_SPPBukti A
            			INNER JOIN Ta_SPP B ON A.No_SPP = B.No_SPP
			INNER JOIN Ta_Pencairan C ON B.No_SPP = C.No_SPP
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE B.Jn_SPP = 'LS'
        AND YEAR(C.Tgl_Cek) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, B.Tgl_SPJ AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_SPJBukti A
            INNER JOIN Ta_SPJ B ON A.No_SPJ = B.No_SPJ
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE YEAR(B.Tgl_SPJ) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (- A.Nilai) Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_TBPRinci A
            			INNER JOIN Ta_TBP B ON A.No_Bukti = B.No_Bukti
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 1) = '5'
          AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_Mutasi A
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 1) = '5'
          AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (A.Debet - A.Kredit) Nilai, B.Tanggal AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_JurnalUmumRinci A
        INNER JOIN Ta_JurnalUmum B ON A.NoBukti = B.NoBukti
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
		WHERE
			LEFT ( A.Kd_Rincian, 1 ) = '5'
          AND YEAR(B.Tanggal) <= @Tahun
			AND B.Posted = 1
    ) AS CombinedTables
    GROUP BY Tahun, Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana
),
`,

	CommonSQL3: `
				CTE_AggregatedData AS (
    SELECT
        Tahun,
        Kd_Desa,
        Kd_Keg,
        COALESCE(Kd_SubRinci, '01') AS Kd_SubRinci,
        Kd_Rincian,
        Sumberdana,
						SUM(CASE WHEN KdPosting = 2 THEN Anggaran ELSE 0 END) AS Anggaran1,
        SUM(CASE WHEN KdPosting = MaxKdPosting THEN (CASE WHEN KdPosting = 2 THEN Anggaran ELSE AnggaranStlhPAK END) ELSE 0 END) AS Anggaran2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 1 THEN Nilai ELSE 0 END) AS Real1,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 2 THEN Nilai ELSE 0 END) AS Real2,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 3 THEN Nilai ELSE 0 END) AS Real3,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 4 THEN Nilai ELSE 0 END) AS Real4,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 5 THEN Nilai ELSE 0 END) AS Real5,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 6 THEN Nilai ELSE 0 END) AS Real6,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 7 THEN Nilai ELSE 0 END) AS Real7,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 8 THEN Nilai ELSE 0 END) AS Real8,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 9 THEN Nilai ELSE 0 END) AS Real9,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 10 THEN Nilai ELSE 0 END) AS Real10,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 11 THEN Nilai ELSE 0 END) AS Real11,
        SUM(CASE WHEN MONTH(Tgl_Bukti) = 12 THEN Nilai ELSE 0 END) AS Real12,
        SUM(Nilai) AS TotalReal
    FROM (
                SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, Anggaran, AnggaranStlhPAK, 0 AS Nilai, NULL AS Tgl_Bukti, A.KdPosting, LKP.MaxKdPosting FROM Ta_AnggaranRinci A
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        LEFT JOIN LatestKdPosting LKP ON A.Tahun = LKP.Tahun AND A.Kd_Desa = LKP.Kd_Desa
        WHERE LEFT(A.Kd_Rincian, 1) = '6'
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, C.Tgl_Cek AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_SPPBukti A
            INNER JOIN Ta_SPP B ON A.No_SPP = B.No_SPP
			INNER JOIN Ta_Pencairan C ON B.No_SPP = C.No_SPP
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE B.Jn_SPP = 'PBY'
          AND YEAR(C.Tgl_Cek) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_TBPRinci A
            			INNER JOIN Ta_TBP B ON A.No_Bukti = B.No_Bukti
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 4) = '6.1.'
          AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (- A.Nilai) Nilai, Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_TBPRinci A
            			INNER JOIN Ta_TBP B ON A.No_Bukti = B.No_Bukti
			LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
        WHERE LEFT(Kd_Rincian, 4) = '6.2.'
        AND YEAR(Tgl_Bukti) <= @Tahun
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (A.Debet - A.Kredit) Nilai, B.Tanggal AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_JurnalUmumRinci A
        INNER JOIN Ta_JurnalUmum B ON A.NoBukti = B.NoBukti
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
		WHERE
			LEFT ( A.Kd_Rincian, 4 ) = '6.2.'
          AND YEAR(B.Tanggal) <= @Tahun
			AND B.Posted = 1
        UNION ALL
        SELECT A.Tahun, A.Kd_Desa, Kd_Keg, '01' AS Kd_SubRinci, Kd_Rincian, Sumberdana, 0 AS Anggaran, 0 AS AnggaranStlhPAK, (A.Kredit - A.Debet) Nilai, B.Tanggal AS Tgl_Bukti, 0 AS KdPosting, 0 AS MaxKdPosting FROM Ta_JurnalUmumRinci A
        INNER JOIN Ta_JurnalUmum B ON A.NoBukti = B.NoBukti
        LEFT JOIN Ta_Pemda tp ON A.Tahun = tp.Tahun
		WHERE
			LEFT ( A.Kd_Rincian, 4 ) = '6.1.'
          AND YEAR(B.Tanggal) <= @Tahun
			AND B.Posted = 1
    ) AS CombinedTables
    GROUP BY Tahun, Kd_Desa, Kd_Keg, Kd_SubRinci, Kd_Rincian, Sumberdana
),
`,

	TaggingSQL: `
			CTE_Tagging AS (
    SELECT
        B.Tahun,
        B.Kd_Desa,
        B.Kd_Keg,
        B.No_ID,
		(SELECT A.ID_Tagging + '; '
        FROM Ta_Tagging AS A
        WHERE A.Tahun = B.Tahun
          AND A.Kd_Desa = B.Kd_Desa
          AND A.Kd_Keg = B.Kd_Keg
          AND A.No_ID = B.No_ID
        FOR XML PATH (''), TYPE).value('.', 'NVARCHAR(MAX)') AS Tagging
    FROM Ta_Tagging B
    GROUP BY B.Tahun, B.Kd_Desa, B.Kd_Keg, B.No_ID
)`,

	ConditionalSQL2024: `
SELECT
    ROW_NUMBER() OVER(ORDER BY (SELECT NULL)) AS ID,
    A.Tahun,
    C.Kd_Prov + C.Kd_Kab + REPLACE(A.Kd_Desa, '.', '') AS Kode_Desa,
    REPLACE(SUBSTRING(A.Kd_Keg, 9, 10), '.', '') AS Kode_Kegiatan,
    SUBSTRING(A.Kd_Keg, 9, 10) AS Kd_Kegiatan,
    A.Kd_SubRinci AS Kode_Paket,
    REPLACE(A.Kd_Rincian, '.', '') AS Kode_Rekening,
    A.Sumberdana AS Kode_Sumber,
    COALESCE(B.Tagging, '') AS Tagging,
    SUM(A.Anggaran1) AS Anggaran1,
    SUM(A.Anggaran2) AS Anggaran2,
    SUM(A.Real1) AS Real1,
    SUM(A.Real2) AS Real2,
    SUM(A.Real3) AS Real3,
    SUM(A.Real4) AS Real4,
    SUM(A.Real5) AS Real5,
    SUM(A.Real6) AS Real6,
    SUM(A.Real7) AS Real7,
    SUM(A.Real8) AS Real8,
    SUM(A.Real9) AS Real9,
    SUM(A.Real10) AS Real10,
    SUM(A.Real11) AS Real11,
    SUM(A.Real12) AS Real12,
    SUM(A.TotalReal) AS TotalReal,
    C.Kd_Prov+C.Kd_Kab AS Kode_Pemda,
    C.Nama_Pemda,
    F.Nama_Obyek AS Nama_Rekening,
    G.Nama_Sumber,
    H.Nama_Desa,
    D.Nama_Kegiatan,
    E.Nama_Paket,
    I.ID_Tipologi
FROM
    CTE_AggregatedData A
LEFT JOIN CTE_Tagging B ON A.Tahun = B.Tahun AND A.Kd_Desa = B.Kd_Desa AND A.Kd_Keg = B.Kd_Keg AND A.Kd_SubRinci = B.No_ID
INNER JOIN Ref_Desa H ON A.Kd_Desa = H.Kd_Desa
INNER JOIN Ta_Desa I ON A.Kd_Desa = I.Kd_Desa
LEFT JOIN Ref_Kegiatan D ON SUBSTRING ( A.Kd_Keg, 9, 9 ) = D.ID_Keg
LEFT JOIN Ta_KegiatanOutput E ON A.Kd_Desa = E.Kd_Desa AND A.Kd_Keg = E.Kd_Keg AND A.Kd_SubRinci = E.No_ID
INNER JOIN Ref_Rek4 F ON A.Kd_Rincian = F.Obyek
INNER JOIN Ref_Sumber G ON A.Sumberdana = G.Kode,
Ta_Pemda C
WHERE I.Tahun = @Tahun
GROUP BY
    A.Tahun,
    C.Kd_Prov,
    C.Kd_Kab,
    C.Nama_Pemda,
    A.Kd_Desa,
    A.Kd_Keg,
    A.Kd_SubRinci,
    A.Kd_Rincian,
    A.Sumberdana,
    B.Tagging,
    H.Nama_Desa,
    I.ID_Tipologi,
    D.Nama_Kegiatan,
    E.Nama_Paket,
    F.Nama_Obyek,
    G.Nama_Sumber
ORDER BY
    A.Tahun,
    C.Kd_Prov,
    C.Kd_Kab,
    A.Kd_Desa,
    A.Kd_Keg,
    A.Kd_SubRinci,
    A.Kd_Rincian,
    A.Sumberdana,
    B.Tagging
`,

	ConditionalSQLBefore2024: `
SELECT
    ROW_NUMBER() OVER(ORDER BY (SELECT NULL)) AS ID,
     A.Tahun,
    C.Kd_Prov + C.Kd_Kab + REPLACE(A.Kd_Desa, '.', '') AS Kode_Desa,
    REPLACE(SUBSTRING(A.Kd_Keg, 9, 10), '.', '') AS Kode_Kegiatan,
    SUBSTRING(A.Kd_Keg, 9, 10) AS Kd_Kegiatan,
    A.Kd_SubRinci AS Kode_Paket,
    REPLACE(A.Kd_Rincian, '.', '') AS Kode_Rekening,
    A.Sumberdana AS Kode_Sumber,
    COALESCE(B.Tagging, '') AS Tagging,
    SUM(A.Anggaran1) AS Anggaran1,
    SUM(A.Anggaran2) AS Anggaran2,
        SUM(A.Real1) AS Real1,
    SUM(A.Real2) AS Real2,
    SUM(A.Real3) AS Real3,
    SUM(A.Real4) AS Real4,
    SUM(A.Real5) AS Real5,
    SUM(A.Real6) AS Real6,
    SUM(A.Real7) AS Real7,
    SUM(A.Real8) AS Real8,
    SUM(A.Real9) AS Real9,
    SUM(A.Real10) AS Real10,
    SUM(A.Real11) AS Real11,
    SUM(A.Real12) AS Real12,
    SUM(A.TotalReal) AS TotalReal,
    C.Kd_Prov+C.Kd_Kab AS Kode_Pemda,
    C.Nama_Pemda,
    F.Nama_Obyek AS Nama_Rekening,
    G.Nama_Sumber,
    H.Nama_Desa,
    D.Nama_Kegiatan,
    NULL AS ID_Tipologi
FROM
    CTE_AggregatedData A
LEFT JOIN CTE_Tagging B ON A.Tahun = B.Tahun AND A.Kd_Desa = B.Kd_Desa AND A.Kd_Keg = B.Kd_Keg AND A.Kd_SubRinci = B.No_ID
INNER JOIN Ref_Desa H ON A.Kd_Desa = H.Kd_Desa
LEFT JOIN Ref_Kegiatan D ON SUBSTRING ( A.Kd_Keg, 9, 9 ) = D.ID_Keg
LEFT JOIN Ta_KegiatanOutput E ON A.Kd_Desa = E.Kd_Desa AND A.Kd_Keg = E.Kd_Keg AND A.Kd_SubRinci = E.No_ID
INNER JOIN Ref_Rek4 F ON A.Kd_Rincian = F.Obyek
INNER JOIN Ref_Sumber G ON A.Sumberdana = G.Kode,
Ta_Pemda C
			WHERE A.Tahun = @Tahun
GROUP BY
    A.Tahun,
    C.Kd_Prov,
    C.Kd_Kab,
    C.Nama_Pemda,
    A.Kd_Desa,
    A.Kd_Keg,
    A.Kd_SubRinci,
    A.Kd_Rincian,
    A.Sumberdana,
    B.Tagging,
    H.Nama_Desa,
    D.Nama_Kegiatan,
    E.Nama_Paket,
    F.Nama_Obyek,
    G.Nama_Sumber
ORDER BY
    A.Tahun,
    C.Kd_Prov,
    C.Kd_Kab,
    A.Kd_Desa,
    A.Kd_Keg,
    A.Kd_SubRinci,
    A.Kd_Rincian,
    A.Sumberdana,
    B.Tagging
`,
}