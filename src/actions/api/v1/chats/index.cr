class Api::V1::Chats::Index < ApiAction
  get "/chats" do
    # TODO(11b): pagination, basic filtering (contentious)
    chats = ChatQuery.new.preload_characters.creator_id(current_user.id)

    serialized_chats = ChatSerializer.for_collection(chats)
    json({chats: serialized_chats})
  end
end
