import java.util.Base64

plugins {
    id("com.android.library")
    id("org.jetbrains.kotlin.android")
    id("maven-publish")
    id("signing")
    id("tech.yanand.maven-central-publish") version "1.2.0"
}

android {
    namespace = "io.github.perawallet"
    publishing {
        singleVariant("release") {
            withSourcesJar()
        }
    }
    compileSdk = libs.versions.android.compileSdk.get().toInt()
}

afterEvaluate {
    val versionTag = System.getenv("VERSION_TAG") ?: "0.1.0"
    publishing {
        publications {
            // Publication for AlgoSDK
            create<MavenPublication>("AlgoSDKRelease") {
                artifact(file("AlgoSDK.aar")) {
                    extension = "aar"
                }
                groupId = "io.github.perawallet"
                artifactId = "algorand-go-mobile-sdk"
                version = versionTag
                setupPom("AlgoSDK")
            }

            // Publication for PeraCompactDecimalFormat
            create<MavenPublication>("PeraCompactDecimalFormatRelease") {
                artifact(file("peracompactdecimalformat.aar")) {
                    extension = "aar"
                }
                groupId = "io.github.perawallet"
                artifactId = "peracompactdecimalformat"
                version = versionTag
                setupPom("PeraCompactDecimalFormat")
            }

            // Publication for WalletConnect
            create<MavenPublication>("PeraWalletConnectRelease") {
                artifact(file("perawalletconnect.aar")) {
                    extension = "aar"
                }
                groupId = "io.github.perawallet"
                artifactId = "perawalletconnect"
                version = versionTag
                setupPom("PeraWalletConnect")
            }
        }

        repositories {
            maven {
                name = "Local"
                url = uri(layout.buildDirectory.dir("repos/bundles").get().asFile.toURI())
            }
        }
    }
}

signing {
    // About GPG signing, please refer to https://central.sonatype.org/publish/requirements/gpg/
    val signingKey = System.getenv("GPG_PRIVATE_KEY") ?: ""
    val signingPassword = System.getenv("GPG_PASSPHRASE") ?: ""
    useInMemoryPgpKeys(signingKey, signingPassword)
    sign(publishing.publications)
}

val username = System.getenv("OSSRH_USERNAME") ?: ""
val password = System.getenv("OSSRH_PASSWORD") ?: ""

mavenCentral {
    repoDir = layout.buildDirectory.dir("repos/bundles")
    // Token for Publisher API calls obtained from Sonatype official,
    // it should be Base64 encoded of "username:password".
    authToken = Base64.getEncoder().encodeToString("$username:$password".toByteArray())
    // Whether the upload should be automatically published or not. Use 'USER_MANAGED' if you wish to do this manually.
    // This property is optional and defaults to 'AUTOMATIC'.
    publishingType = "AUTOMATIC"
    // Max wait time for status API to get 'PUBLISHING' or 'PUBLISHED' status when the publishing type is 'AUTOMATIC',
    // or additionally 'VALIDATED' when the publishing type is 'USER_MANAGED'.
    // This property is optional and defaults to 60 seconds.
    maxWait = 500
}

tasks.register("publishAllToMavenLocal") {
    dependsOn("publishAlgoSDKReleasePublicationToMavenLocal")
    dependsOn("publishPeraCompactDecimalFormatReleasePublicationToMavenLocal")
    dependsOn("publishWalletConnectReleasePublicationToMavenLocal")
}

// Helper function to configure POM metadata
fun MavenPublication.setupPom(libName: String) {
    pom {
        packaging = "aar"
        this.name.set(libName)
        this.description.set("$libName: Android Library for Pera Wallet")
        this.url.set("https://github.com/perawallet/algorand-go-mobile-sdk.git")
        this.inceptionYear.set("2025")

        licenses {
            license {
                this.name.set("The Apache License, Version 2.0")
                this.url.set("https://github.com/perawallet/algorand-go-mobile-sdk/blob/main/LICENSE")
            }
        }

        developers {
            developer {
                this.id.set("PeraWallet")
                this.name.set("Pera Wallet")
                this.email.set("contact@perawallet.app")
            }
        }

        scm {
            this.connection.set("scm:git@github.com:perawallet/algorand-go-mobile-sdk.git")
            this.developerConnection.set("scm:git@github.com:perawallet/algorand-go-mobile-sdk.git")
            this.url.set("https://github.com/perawallet/algorand-go-mobile-sdk.git")
        }
    }
}
