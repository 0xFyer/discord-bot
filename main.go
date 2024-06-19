package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/0xfyer/discord-bot/state"
	"github.com/bwmarrin/discordgo"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_SECRET"))
	if err != nil {
		return err
	}

	// When the bot connects to a server
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println(r.Guilds)
	})

	// When someone uses a slash command (i.e. /blackjack)
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "blackjack":
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":white_check_mark:",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			ch, err := s.State.Channel(i.ChannelID)
			if err != nil {
				return
			}

			if ch.IsThread() {
				return
			}

			if state.Info.ParentChannelHasGame(i.ChannelID) {
				if !state.Info.GameHasPlayer(i.ChannelID, i.Member.User.ID) {
					state.Info.AddPlayer(i.ChannelID, i.Member.User.ID)
				}
				return
			}

			thread, err := s.ThreadStart(i.ChannelID, "Blackjack Table", 0, 60)
			if err != nil {
				fmt.Println(err)
				return
			}

			header, err := s.ChannelMessageSend(thread.ID, fmt.Sprintf("<@%s>", i.Member.User.ID))
			if err != nil {
				return
			}

			state.Info.AddNewGame(i.ChannelID, thread.ID, header.ID, i.Member.User.ID)

			_, err = s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
				Flags:   1 << 6,
				Content: "Your move...",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Hit",
								Emoji: &discordgo.ComponentEmoji{
									Name: "💥",
								},
								Style:    discordgo.PrimaryButton,
								CustomID: "hit",
								Disabled: false,
							},
							discordgo.Button{
								Label: "Stand",
								Emoji: &discordgo.ComponentEmoji{
									Name: "🤌🏻",
								},
								Style:    discordgo.SecondaryButton,
								CustomID: "stand",
								Disabled: false,
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
	// When the bot views a message
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		args := strings.Split(m.Message.Content, " ")

		switch args[0] {
		case "/blackjack":
			ch, err := s.State.Channel(m.ChannelID)
			if err != nil {
				return
			}

			if ch.IsThread() {
				return
			}

			if state.Info.ParentChannelHasGame(m.ChannelID) {
				if !state.Info.GameHasPlayer(m.ChannelID, m.Author.ID) {
					state.Info.AddPlayer(m.ChannelID, m.Author.ID)
				}
				return
			}

			thread, err := s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
				Name:                "Blackjack Table",
				AutoArchiveDuration: 60,
				Invitable:           true,
				RateLimitPerUser:    10,
			})

			if err != nil {
				return
			}

			header, err := s.ChannelMessageSend(thread.ID, "Greetings")
			if err != nil {
				return
			}

			state.Info.AddNewGame(m.ChannelID, thread.ID, header.ID, m.Author.ID)

			header, err = s.ChannelMessageEdit(thread.ID, header.ID, "Changed")
			if err != nil {
				return
			}

			_, err = s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
				Flags:   1 << 6,
				Content: "Your move...",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Hit",
								Emoji: &discordgo.ComponentEmoji{
									Name: "💥",
								},
								Style:    discordgo.PrimaryButton,
								CustomID: "hit",
								Disabled: false,
							},
							discordgo.Button{
								Label: "Stand",
								Emoji: &discordgo.ComponentEmoji{
									Name: "🤌🏻",
								},
								Style:    discordgo.SecondaryButton,
								CustomID: "stand",
								Disabled: false,
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	// Handle Component Interactions
	session.AddHandler(func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		switch m.Type {
		case discordgo.InteractionMessageComponent:
			if m.MessageComponentData().CustomID == "hit" {

			}

			if m.MessageComponentData().CustomID == "stand" {
			}

			err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "Waiting on the table...",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label: "Hit",
									Emoji: &discordgo.ComponentEmoji{
										Name: "💥",
									},
									Style:    discordgo.PrimaryButton,
									CustomID: "hit",
									Disabled: true,
								},
								discordgo.Button{
									Label: "Stand",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🤌🏻",
									},
									Style:    discordgo.SecondaryButton,
									CustomID: "stand",
									Disabled: true,
								},
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	})

	err = session.Open()
	if err != nil {
		return err
	}
	defer session.Close()

	// When someone types /blackjack
	bj_cmd := &discordgo.ApplicationCommand{
		Name:        "blackjack",
		Description: "Launch or join a game of blackjack in the current channel",
	}

	_, err = session.ApplicationCommandCreate(session.State.Application.ID, session.State.Application.GuildID, bj_cmd)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	return nil
}
