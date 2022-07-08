package util

//func GenerateJwt(u *database.User) (string, error) {
//	expiredTime := time.Now().Add(DefaultValidityPeriod)
//	stdClaims := jwt.StandardClaims{
//		ExpiresAt: expiredTime.Unix(),
//		IssuedAt:  time.Now().Unix(),
//		Id:        string(u.ID),
//		Issuer:    AppIssuer,
//	}
//	uClaims := database.UserStdClaims{StandardClaims: stdClaims, User: u}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
//	tokenString, err := token.SignedString(AppSecret)
//	if err != nil {
//		Logger.WithFields(logrus.Fields{
//			"err":  err,
//			"user": u,
//		}).Error("generate jwt error")
//		return "", err
//	}
//	return tokenString, nil
//}

//func ParseJwt(c *gin.Context) (*database.User, error) {
//	tokenString, err := GetToken(c)
//	if err != nil {
//		return nil, err
//	}
//	if tokenString == "" {
//		Logger.Debug("tokenString not allow nil")
//		return nil, TokenEmptyErrorTemplate
//	}
//	claims := database.UserStdClaims{}
//	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return []byte(AppSecret), nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	return claims.User, err
//}
