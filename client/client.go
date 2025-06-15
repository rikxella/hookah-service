package main

import (
	"context"
	"flag"
	"log"
	"os"

	pb "hookah-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	brandPrefix := flag.String("brand", "", "Префикс для поиска бренда (например: 'al')")
	flavorBrand := flag.String("flavor-brand", "", "Бренд для поиска вкуса (например: 'Al Fakher')")
	flavorPrefix := flag.String("flavor", "", "Префикс для поиска вкуса (например: 'apple')")
	flag.Parse()

	if *brandPrefix == "" && (*flavorBrand == "" || *flavorPrefix == "") {
		log.Println("Использование:")
		log.Println("  Поиск брендов: client -brand <префикс>")
		log.Println("  Поиск вкусов: client -flavor-brand <бренд> -flavor <префикс>")
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := pb.NewTobaccoSearchServiceClient(conn)
	ctx := context.Background()

	if *brandPrefix != "" {
		brandResp, err := client.BrandTobacco(ctx, &pb.BrandTobaccoSearch{
			BrandPrefix: *brandPrefix,
		})
		if err != nil {
			log.Fatalf("Ошибка поиска брендов: %v", err)
		}
		log.Printf("Найдены бренды по запросу '%s':", *brandPrefix)
		for _, brand := range brandResp.Results {
			log.Printf("  • %s", brand)
		}
	}

	if *flavorBrand != "" && *flavorPrefix != "" {
		nameResp, err := client.NameTobacco(ctx, &pb.NameTobaccoSearch{
			Brand:      *flavorBrand,
			NamePrefix: *flavorPrefix,
		})
		if err != nil {
			log.Fatalf("Ошибка поиска вкусов: %v", err)
		}
		log.Printf("Найдены вкусы бренда '%s' по запросу '%s':", *flavorBrand, *flavorPrefix)
		for _, flavor := range nameResp.Results {
			log.Printf("  • %s", flavor)
		}
	}
}
