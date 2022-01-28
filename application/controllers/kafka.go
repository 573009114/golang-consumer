package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

//消费者
func Consumer() {
	var cfg cfg

	if _, err := toml.DecodeFile("./config/config.toml", &cfg); err != nil {
		fmt.Println(err)
	}
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	c, err := cluster.NewConsumer(strings.Split(cfg.Kfk.Endpoint, ","), cfg.Kfk.GroupID, strings.Split(cfg.Kfk.Topic, ","), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		notify := c.Notifications()
		for {
			select {
			case err := <-errors:
				fmt.Println(err)
			case <-notify:

			}
		}
	}(c)
	for msg := range c.Messages() {
		// fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		c.MarkOffset(msg, "")
		BulkIndexToEs(cfg.Kfk.Topic, msg.Value)
	}

}
