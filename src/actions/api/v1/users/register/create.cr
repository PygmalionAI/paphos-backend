class Api::V1::Users::Register::Create < ApiAction
  include Api::Auth::SkipRequireAuthToken

  post "/users/register" do
    user = RegisterUser.create!(params)
    json({token: UserToken.generate(user)})
  end
end
