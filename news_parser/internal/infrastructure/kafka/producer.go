package main

//type KafkaProducer struct {
//	producer *kafka.Producer
//}
//
//func NewKafkaProducer(servers []string, producerName string) (*KafkaProducer, error) {
//	// Настройка конфигурации продюсера нужно автоматизировать
//	config := &kafka.ConfigMap{
//		"bootstrap.servers":  strings.Join(servers, ","), // Адреса серверов в кластере Kafka. Поскольку все брокеры расположены на одном хосте, то у них один адрес, но слушают они на разных портах.
//		"acks":               "all",                      // Ожидание подтверждения от всех реплик
//		"enable.idempotence": true,                       // включить идемпотентность (без дублей)
//		"retries":            3,                          // количество ретраев при ошибке
//		"retry.backoff.ms":   100,                        // пауза между ретраями
//		"client.id":          producerName,               // имя продюсера
//	}
//	producer, err := kafka.NewProducer(config)
//	if err != nil {
//		fmt.Printf("Failed to create producer: %s\n", err)
//		os.Exit(1)
//		return nil, err
//	}
//	return &KafkaProducer{producer}, nil
//}
//
//func (k *KafkaProducer) Close() {
//	k.producer.Close()
//}
//
//func (k *KafkaProducer) Produce(msg *dto.ArticleMessage, topic string) error {
//	msgJson, err := json.Marshal(msg)
//	if err != nil {
//		return err
//	}
//	deliveryChan := make(chan kafka.Event)
//	err = k.producer.Produce(&kafka.Message{
//		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
//		Value:          msgJson,
//	}, deliveryChan)
//	if err != nil {
//		return err
//	}
//	result := <-deliveryChan
//	kafkaMessage := result.(*kafka.Message)
//	if kafkaMessage.TopicPartition.Error != nil {
//		fmt.Printf("Ошибка доставки: %v\n", kafkaMessage.TopicPartition.Error)
//	} else {
//		fmt.Printf("Сообщение доставлено в %v [%d] offset %v\n",
//			*kafkaMessage.TopicPartition.Topic, kafkaMessage.TopicPartition.Partition, kafkaMessage.TopicPartition.Offset)
//	}
//	close(deliveryChan)
//	return nil
//}
//
//func consumer() {
//	// Настройка конфигурации консьюмера
//	config := &kafka.ConfigMap{
//		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094", //адреса Kafka-брокеров;
//		"group.id":          "myGroup",                                      //идентификатор группы потребителей
//		"auto.offset.reset": "smallest",                                     //стратегия поведения при отсутствии смещения, например, smallest — для чтения с начала.
//	}
//
//	// Инициализация консьюмера
//	consumer, err := kafka.NewConsumer(config)
//	if err != nil {
//		panic(fmt.Sprintf("Failed to create consumer: %v", err))
//	}
//
//	err = consumer.SubscribeTopics([]string{"async-topic"}, nil) // Указываем подписку на топики async-topic и sync-topic.
//
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Consumer initialized")
//	run := true
//	//consumer.
//	for run {
//		//ev := consumer.Poll(100)
//		ev, err := consumer.ReadMessage(-1)
//		if err != nil {
//			fmt.Println("Consumer error:", err)
//		}
//		fmt.Printf("Message: %v", *ev)
//		//fmt.Println("Message:", string(ev.Value))
//		//switch e := ev.(type) {
//		//case *queue.Message:
//		//	// application-specific processing
//		//case queue.Error:
//		//	fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
//		//	run = false
//		//default:
//		//	fmt.Printf("Ignored %v\n", e)
//		//	time.Sleep(5 * time.Second)
//		//}
//	}
//
//	consumer.Close()
//
//}

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"sync"
)

func producer() {
	// Настройка конфигурации продюсера
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094", // Адреса серверов в тестовом кластере Kafka. Поскольку все брокеры расположены на одном хосте, то у них один адрес, но слушают они на разных портах.
		"acks":              "all",                                          // Ожидание подтверждения от всех реплик
		"client.id":         "myProducer",                                   //имя продюсера
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	defer producer.Close()
	fmt.Println("Producer initialized")
	topic := "async-topic"

	for _, word := range []string{"this", "is", "asynchronous", "message", "delivery", "in", "kafka", "with", "Go", "Client"} {
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
	fmt.Println("Сообщения отправлены")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Асинхронная обработка событий продюсера
	go func() {
		defer wg.Done()
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v, msg %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset,
						string(ev.Value))
				}
			}
		}
	}()
	wg.Wait()
}

func consumer() {
	// Настройка конфигурации консьюмера
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094", //адреса Kafka-брокеров;
		"group.id":          "myGroup",                                      //идентификатор группы потребителей
		"auto.offset.reset": "smallest",                                     //стратегия поведения при отсутствии смещения, например, smallest — для чтения с начала.
	}

	// Инициализация консьюмера
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %v", err))
	}

	err = consumer.SubscribeTopics([]string{"async-topic"}, nil) // Указываем подписку на топики async-topic и sync-topic.

	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer initialized")
	run := true
	//consumer.
	for run {
		//ev := consumer.Poll(100)
		ev, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println("Consumer error:", err)
		}
		fmt.Printf("Message: %v", string(ev.Value))
		//fmt.Println("Message:", string(ev.Value))
		//switch e := ev.(type) {
		//case *kafka.Message:
		//	// application-specific processing
		//case kafka.Error:
		//	fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
		//	run = false
		//default:
		//	fmt.Printf("Ignored %v\n", e)
		//	time.Sleep(5 * time.Second)
		//}
	}

	consumer.Close()

}

func main() {
	go producer()
	consumer()

}
