package seed

import (
	"encoding/json"
	"fmt"
	"fullstack-journal/app/entity"
	"io"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return json.Unmarshal(body, target)
}

func IndonesiaLocationSeed(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// Seed Province
		var provinceList []map[string]interface{}
		err := fetchJSON("https://ibnux.github.io/data-indonesia/propinsi.json", &provinceList)
		if err != nil {
			return fmt.Errorf("gagal fetch provinsi: %w", err)
		}

		for _, item := range provinceList {
			id, err := strconv.Atoi(item["id"].(string))
			if err != nil {
				return fmt.Errorf("gagal konversi Province ID: %s", item["id"])
			}

			prov := entity.Province{
				ID:       int64(id),
				Name:     item["nama"].(string),
				Latitude: item["latitude"].(float64),
				Longitude: item["longitude"].(float64),
			}

			if err := tx.FirstOrCreate(&prov, entity.Province{ID: prov.ID}).Error; err != nil {
				return fmt.Errorf("gagal insert province %s: %w", prov.Name, err)
			}
		}

		// Seed Regency
		for _, prov := range provinceList {
			provID, err := strconv.Atoi(prov["id"].(string))
			if err != nil {
				return fmt.Errorf("gagal konversi Province ID: %s", prov["id"])
			}

			var regencyList []map[string]interface{}
			url := fmt.Sprintf("https://ibnux.github.io/data-indonesia/kabupaten/%s.json", prov["id"])
			err = fetchJSON(url, &regencyList)
			if err != nil {
				return fmt.Errorf("gagal fetch kabupaten untuk provinsi %s: %w", prov["id"], err)
			}

			for _, regency := range regencyList {
				id, err := strconv.Atoi(regency["id"].(string))
				if err != nil {
					return fmt.Errorf("gagal konversi Regency ID: %s", regency["id"])
				}

				reg := entity.Regency{
					ID:         int64(id),
					ProvinceId: int64(provID),
					Name:       regency["nama"].(string),
					Latitude:   regency["latitude"].(float64),
					Longitude:  regency["longitude"].(float64),
				}

				if err := tx.FirstOrCreate(&reg, entity.Regency{ID: reg.ID}).Error; err != nil {
					return fmt.Errorf("gagal insert regency %s: %w", reg.Name, err)
				}
			}
		}

		// Seed SubDistrict
		var regencyList []entity.Regency
		if err := tx.Find(&regencyList).Error; err != nil {
			return fmt.Errorf("gagal ambil data regency dari DB: %w", err)
		}

		for _, regency := range regencyList {
			var subDistrictList []map[string]interface{}
			url := fmt.Sprintf("https://ibnux.github.io/data-indonesia/kecamatan/%d.json", regency.ID)
			err := fetchJSON(url, &subDistrictList)
			if err != nil {
				return fmt.Errorf("gagal fetch kecamatan untuk regency %s: %w", regency.Name, err)
			}

			for _, subDistrict := range subDistrictList {
				id, err := strconv.Atoi(subDistrict["id"].(string))
				if err != nil {
					return fmt.Errorf("gagal konversi SubDistrict ID: %s", subDistrict["id"])
				}

				subDst := entity.SubDistrict{
					ID:        int64(id),
					RegencyId: regency.ID,
					Name:      subDistrict["nama"].(string),
					Latitude:  subDistrict["latitude"].(float64),
					Longitude: subDistrict["longitude"].(float64),
				}

				if err := tx.FirstOrCreate(&subDst, entity.SubDistrict{ID: subDst.ID}).Error; err != nil {
					return fmt.Errorf("gagal insert subdistrict %s: %w", subDst.Name, err)
				}
			}
		}

		// Semua berhasil, commit
		return nil
	})

	if err != nil {
		log.Println("Gagal melakukan seeding data lokasi:", err)
	} else {
		log.Println("Sukses melakukan seeding semua data lokasi (provinsi, kabupaten, kecamatan).")
	}
}


	
