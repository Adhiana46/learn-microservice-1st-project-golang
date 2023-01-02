version := "latest"

build-images:
	echo "Building images version: $(version)"
	docker build -t adhiana46/ms-golang-client-react:$(version) ./client-react/
	docker build -t adhiana46/ms-golang-post-service:$(version) ./posts/
	docker build -t adhiana46/ms-golang-comment-service:$(version) ./comments/
	docker build -t adhiana46/ms-golang-event-bus:$(version) ./event-bus/
	docker build -t adhiana46/ms-golang-query-service:$(version) ./query/
	# docker build -t adhiana46/ms-golang-moderation-service:$(version) ./moderation/
	echo "Done building images."
push-images:
	echo "Pushing images version $(version) to docker hub"
	docker push adhiana46/ms-golang-client-react:$(version)
	docker push adhiana46/ms-golang-post-service:$(version)
	docker push adhiana46/ms-golang-comment-service:$(version)
	docker push adhiana46/ms-golang-event-bus:$(version)
	docker push adhiana46/ms-golang-query-service:$(version)
	# docker push adhiana46/ms-golang-moderation-service:$(version)
	echo "Done pushing images."
k8s-apply:
	kubectl apply -f ./infra/k8s
k8s-delete:
	kubectl delete -f ./infra/k8s
k8s-rollout-restart:
	kubectl rollout restart deployment/client-react
	kubectl rollout restart deployment/comments
	kubectl rollout restart deployment/event-bus
	kubectl rollout restart deployment/moderation
	kubectl rollout restart deployment/posts
	kubectl rollout restart deployment/query
minikube-clean:
	minikube image rm \
	docker.io/adhiana46/ms-golang-client-react \
	docker.io/adhiana46/ms-golang-post-service \
	docker.io/adhiana46/ms-golang-comment-service \
	docker.io/adhiana46/ms-golang-event-bus \
	docker.io/adhiana46/ms-golang-query-service
	# docker.io/adhiana46/ms-golang-moderation-service \