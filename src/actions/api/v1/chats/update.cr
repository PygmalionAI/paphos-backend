class Api::V1::Chats::Update < ApiAction
  include CheckCurrentUser

  patch "/chats/:chat_id" do
    chat = ChatQuery.new.find(chat_id)
    ensure_owned_by_current_user!(chat)

    updated_chat = UpdateChat.update!(chat, params, current_user: current_user)

    json({chat: ChatSerializer.new(ChatQuery.preload_characters(updated_chat))})
  end
end
