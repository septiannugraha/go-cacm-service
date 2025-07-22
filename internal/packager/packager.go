package packager

import (
	"fmt"
	"os"
	"time"

	"github.com/septiannugraha/go-cacm-service/internal/models"
	"github.com/septiannugraha/go-cacm-service/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Packager handles the conversion of data to protobuf and serialization
type Packager struct{}

// NewPackager creates a new packager instance
func NewPackager() *Packager {
	return &Packager{}
}

// PackageQueryResults converts CacheData to protobuf and saves to file
func (p *Packager) PackageQueryResults(data []models.CacheData, filePath string) error {
	results := &pb.QueryResults{
		Rows: make([]*pb.QueryResultItem, 0, len(data)),
	}

	for _, item := range data {
		pbItem := &pb.QueryResultItem{
			Id:             int32(item.ID),
			Tahun:          item.Tahun,
			KodeDesa:       item.KodeDesa,
			KodeKegiatan:   item.KodeKegiatan,
			KodePaket:      item.KodePaket,
			KodeRekening:   item.KodeRekening,
			KodeSumber:     item.KodeSumber,
			Tagging:        item.Tagging,
			Anggaran1:      int64(item.Anggaran1),
			Anggaran2:      int64(item.Anggaran2),
			Real1:          int64(item.Real1),
			Real2:          int64(item.Real2),
			Real3:          int64(item.Real3),
			Real4:          int64(item.Real4),
			Real5:          int64(item.Real5),
			Real6:          int64(item.Real6),
			Real7:          int64(item.Real7),
			Real8:          int64(item.Real8),
			Real9:          int64(item.Real9),
			Real10:         int64(item.Real10),
			Real11:         int64(item.Real11),
			Real12:         int64(item.Real12),
			Totalreal:      int64(item.TotalReal),
			KodePemda:      item.KodePemda,
			NamaPemda:      item.NamaPemda,
			NamaRekening:   item.NamaRekening,
			NamaSumber:     item.NamaSumber,
			NamaDesa:       item.NamaDesa,
			NamaKegiatan:   item.NamaKegiatan,
			CreatedAt:      item.CreatedAt.Unix(),
			UploadId:       int32(time.Now().Unix() % 2147483647), // Convert to int32
		}

		// Handle nullable fields
		if item.NamaPaket != nil {
			pbItem.NamaPaket = wrapperspb.String(*item.NamaPaket)
		}
		if item.IDTipologi != nil {
			pbItem.IdTipologi = wrapperspb.String(*item.IDTipologi)
		}

		results.Rows = append(results.Rows, pbItem)
	}

	return p.saveToFile(results, filePath)
}

// PackageBelanjaPerBidangPerSumber packages expense by field data
func (p *Packager) PackageBelanjaPerBidangPerSumber(data []models.BelanjaPerBidangPerSumber, filePath string) error {
	results := &pb.BelanjaPerBidangPerSumberResults{
		Rows: make([]*pb.BelanjaPerBidangPerSumberResult, 0, len(data)),
	}

	for _, item := range data {
		pbItem := &pb.BelanjaPerBidangPerSumberResult{
			Tahun:       fmt.Sprintf("%d", item.Tahun),
			KodeProv:    item.KodeProv,
			NamaProv:    item.NamaProv,
			KodePemda:   item.KodePemda,
			NamaPemda:   item.NamaPemda,
			KodeKec:     item.KodeKec,
			NamaKec:     item.NamaKec,
			Kodedesa:    item.KodeDesa,
			NamaDesa:    item.NamaDesa,
			KodePosting: item.KodePosting,
			SumberDana:  item.SumberDana,
			AnggBid01:   item.AnggBid01,
			RealBid01:   item.RealBid01,
			AnggBid02:   item.AnggBid02,
			RealBid02:   item.RealBid02,
			AnggBid03:   item.AnggBid03,
			RealBid03:   item.RealBid03,
			AnggBid04:   item.AnggBid04,
			RealBid04:   item.RealBid04,
			AnggBid05:   item.AnggBid05,
			RealBid05:   item.RealBid05,
			Currentdate: item.CurrentDate.Format("2006-01-02"),
		}

		results.Rows = append(results.Rows, pbItem)
	}

	return p.saveToFile(results, filePath)
}

// saveToFile serializes protobuf message to binary file
func (p *Packager) saveToFile(message proto.Message, filePath string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// generateUploadID generates a unique upload ID
func generateUploadID() string {
	return fmt.Sprintf("upload_%d", time.Now().Unix())
}