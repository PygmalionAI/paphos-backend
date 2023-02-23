class Api::V1::Chats::Create < ApiAction
  post "/chats" do
    chat = SaveChat.create!(params, current_user: current_user)
    json({chat: ChatSerializer.new(chat)})
  end
end
